package glog

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

const (
	backupTimeFormat = "2006-01-02 15-04-05.000"
	compressSuffix   = ".gz"
	defaultMaxSize   = 100
)

var (
	currentTime = time.Now
	osStat      = os.Stat
	megabyte    = 1024 * 1024
)

// 确保我们始终实现 io.WriteCloser
var _ io.WriteCloser = (*LogFile)(nil)

// LogFile 如果MaxBackups和MaxAge都为0，则不会删除旧的日志文件
type LogFile struct {
	Filename   string
	MaxSize    int  // 日志文件获取之前的最大大小(以兆字节为单位)，默认为100兆字节
	MaxAge     int  // 根据旧日志文件保留的最大天数
	MaxBackups int  // 保留的旧日志文件的最大数量，默认是保留所有旧的日志文件
	Compress   bool // 压缩确定是否应压缩轮转的日志文件，使用gzip，默认不进行压缩
	size       int64
	file       *os.File
	mu         sync.Mutex
	millCh     chan bool
	startMill  sync.Once
}

func (lf *LogFile) Write(p []byte) (n int, err error) {
	lf.mu.Lock()
	defer lf.mu.Unlock()

	writeLen := int64(len(p))
	if writeLen > lf.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, lf.max(),
		)
	}

	if lf.file == nil {
		if err = lf.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	if lf.size+writeLen > lf.max() {
		if err := lf.rotate(); err != nil {
			return 0, err
		}
	}

	n, err = lf.file.Write(p)
	lf.size += int64(n)

	return n, err
}

func (lf *LogFile) Close() error {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	return lf.close()
}

func (lf *LogFile) close() error {
	if lf.file == nil {
		return nil
	}
	err := lf.file.Close()
	lf.file = nil
	return err
}

func (lf *LogFile) Rotate() error {
	lf.mu.Lock()
	defer lf.mu.Unlock()
	return lf.rotate()
}

func (lf *LogFile) rotate() error {
	if err := lf.close(); err != nil {
		return err
	}
	if err := lf.openNew(); err != nil {
		return err
	}
	lf.mill()
	return nil
}

func (lf *LogFile) openNew() error {
	err := os.MkdirAll(lf.dir(), 0755)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := lf.filename()
	mode := os.FileMode(0600)
	info, err := osStat(name)
	if err == nil {
		mode = info.Mode()
		newName := backupName(name)
		if err := os.Rename(name, newName); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}

		// 除了linux之外，这在任何地方都是禁止操作的
		if err := chown(name, info); err != nil {
			return err
		}
	}

	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	lf.file = f
	lf.size = 0
	return nil
}

func backupName(name string) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	timestamp := currentTime().Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, timestamp, ext))
}

func (lf *LogFile) openExistingOrNew(writeLen int) error {
	lf.mill()

	filename := lf.filename()
	info, err := osStat(filename)
	if os.IsNotExist(err) {
		return lf.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= lf.max() {
		return lf.rotate()
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return lf.openNew()
	}
	lf.file = file
	lf.size = info.Size()
	return nil
}

func (lf *LogFile) filename() string {
	if lf.Filename != "" {
		return lf.Filename
	}
	name := filepath.Base(os.Args[0]) + "-backup.log"
	return filepath.Join(os.TempDir(), name)
}

func (lf *LogFile) millRunOnce() error {
	if lf.MaxBackups == 0 && lf.MaxAge == 0 && !lf.Compress {
		return nil
	}

	files, err := lf.oldLogFiles()
	if err != nil {
		return err
	}

	var compress, remove []logInfo

	if lf.MaxBackups > 0 && lf.MaxBackups < len(files) {
		preserved := make(map[string]bool)
		var remaining []logInfo
		for _, f := range files {
			// 只统计未压缩的日志文件或压缩日志文件
			fn := f.Name()
			if strings.HasSuffix(fn, compressSuffix) {
				fn = fn[:len(fn)-len(compressSuffix)]
			}
			preserved[fn] = true

			if len(preserved) > lf.MaxBackups {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}
	if lf.MaxAge > 0 {
		diff := time.Duration(int64(24*time.Hour) * int64(lf.MaxAge))
		cutoff := currentTime().Add(-1 * diff)

		var remaining []logInfo
		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}

	if lf.Compress {
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), compressSuffix) {
				compress = append(compress, f)
			}
		}
	}

	for _, f := range remove {
		errRemove := os.Remove(filepath.Join(lf.dir(), f.Name()))
		if err == nil && errRemove != nil {
			err = errRemove
		}
	}
	for _, f := range compress {
		fn := filepath.Join(lf.dir(), f.Name())
		errCompress := compressLogFile(fn, fn+compressSuffix)
		if err == nil && errCompress != nil {
			err = errCompress
		}
	}

	return err
}

func (lf *LogFile) millRun() {
	for range lf.millCh {
		_ = lf.millRunOnce()
	}
}

// 执行旋转后压缩并删除过时的日志文件
func (lf *LogFile) mill() {
	lf.startMill.Do(func() {
		lf.millCh = make(chan bool, 1)
		go lf.millRun()
	})
	select {
	case lf.millCh <- true:
	default:
	}
}

// 返回存储在同一目录中的备份日志文件列表，目录为当前日志文件，按ModTime排序
func (lf *LogFile) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(lf.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	var logFiles []logInfo

	prefix, ext := lf.prefixAndExt()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if t, err := lf.timeFromName(f.Name(), prefix, ext); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		if t, err := lf.timeFromName(f.Name(), prefix, ext+compressSuffix); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
	}

	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

func (lf *LogFile) timeFromName(filename, prefix, ext string) (time.Time, error) {
	if !strings.HasPrefix(filename, prefix) {
		return time.Time{}, errors.New("mismatched prefix")
	}
	if !strings.HasSuffix(filename, ext) {
		return time.Time{}, errors.New("mismatched extension")
	}
	ts := filename[len(prefix) : len(filename)-len(ext)]
	return time.Parse(backupTimeFormat, ts)
}

func (lf *LogFile) max() int64 {
	if lf.MaxSize == 0 {
		return int64(defaultMaxSize * megabyte)
	}
	return int64(lf.MaxSize) * int64(megabyte)
}

func (lf *LogFile) dir() string {
	return filepath.Dir(lf.filename())
}

func (lf *LogFile) prefixAndExt() (prefix, ext string) {
	filename := filepath.Base(lf.filename())
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)] + "-"
	return prefix, ext
}

func compressLogFile(src, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer f.Close()

	fi, err := osStat(src)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %v", err)
	}

	if err := chown(dst, fi); err != nil {
		return fmt.Errorf("failed to chown compressed log file: %v", err)
	}

	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to open compressed log file: %v", err)
	}
	defer gzf.Close()

	gz := gzip.NewWriter(gzf)

	defer func() {
		if err != nil {
			_ = os.Remove(dst)
			err = fmt.Errorf("failed to compress log file: %v", err)
		}
	}()

	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := gzf.Close(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Remove(src); err != nil {
		return err
	}

	return nil
}

//func chown(_ string, _ os.FileInfo) error {
//	return nil
//}

var osChown = os.Chown

func chown(name string, info os.FileInfo) error {
	//f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	//if err != nil {
	//	return err
	//}
	//f.Close()
	//stat := info.Sys().(*syscall.Stat_t)
	//return osChown(name, int(stat.Uid), int(stat.Gid))
	return nil
}

type logInfo struct {
	timestamp time.Time
	os.FileInfo
}

type byFormatTime []logInfo

func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b byFormatTime) Len() int {
	return len(b)
}
