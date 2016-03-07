package vm
import "os"


func (out)File(v VM) {
	v.Attr()["file-pool"] = newFilePool()
	v.SetOut(48, HANDLE_ALL, outFileFunc)
	v.SetOut(49, HANDLE_ALL, outFileFunc)
	v.SetOut(50, 16, outFileFunc)
	v.SetOut(50, 17, outFileFunc)
	v.SetOut(50, 18, outFileFunc)
	v.SetOut(51, 16, outFileFunc)
	v.SetOut(51, 17, outFileFunc)
	v.SetOut(51, 18, outFileFunc)
	v.SetOut(52, HANDLE_ALL, outFileFunc)
	v.SetOut(53, HANDLE_ALL, outFileFunc)
	v.SetOut(54, HANDLE_ALL, outFileFunc)
	v.SetOut(55, HANDLE_ALL, outFileFunc)
}

//48 | 打开文件 | 0 | r0:打开方式<br>r1:文件号<br>r3:文件名字符串 | 打开方式目前只能为1
//49 | 关闭文件 | 文件号 |  |
//50	| 从文件读取数据
//- |  | 16:读取整数 | r1:文件号<br>r2:位置偏移量 | r3的值变为读取的整数
//-	|  | 17:读取浮点数 | r1:文件号<br>r2:位置偏移量 | r3的值变为读取的浮点数
//-	|  | 18:读取字符串 | r1:文件号<br>r2:位置偏移量<br>r3:目标字符串句柄 | r3所指字符串的内容变为读取的字符串
//51	| 向文件写入数据 |
//- |  | 16:写入整数 | r1:文件号<br>r2:位置偏移量<br>r3:整数 |
//-	|  | 17:写入浮点数 | r1:文件号<br>r2:位置偏移量<br>r3:浮点数 |
//-	|  | 18:写入字符串 | r1:文件号<br>r2:位置偏移量<br>r3:字符串 |
//52 | 判断文件位置指针是否指向文件尾 | 0 | r3:文件号 |  Eof
//53 | 获取文件长度 | 0 | r3:文件号 |  Lof
//54 | 获取文件位置指针的位置 | 0 | r3:文件号 |  Loc
//55 | 定位文件位置指针 | 0 | r2:文件号<br>r3:目标位置 |

// 负责文件打开关闭
func outFileFunc(i *Inst) {
	v, p, o := i.VM, i.A.Get(), i.B // port and param
	r0, r1, r2, r3 := &v.r0, &v.r1, &v.r2, &v.r3
	fp := v.attr["file-pool"].(ResPool)
	switch p{
	case 48:
		mode, fd, fn := r0.Get(), r1.Get(), r3.Str()
		log.Debug("Open file '%s' as #%d mode %d", fn, fd, mode)

		res := fp.Get(fd)
		if res == nil {
			log.Error("Open file number %d invalid", fd)
			return
		}
		f, _ := res.Get().(*os.File)

		if f != nil {
			log.Warning("Last file '%s' not closed, will close now", f.Name())
			err := f.Close()
			if err != nil {
				log.Error("Close file '%s' as #%d with error:%s", f.Name(), fd, err)
			}
		}
		res.Set(nil)

		f, err := os.OpenFile(fn, os.O_RDWR | os.O_CREATE, os.ModePerm)
		if err != nil {
			log.Error("Open file '%s' faied: %s", fn, err.Error())
		}else {
			res.Set(f)
		}
	case 49:
		fd := o.Get()

		res := fp.Get(fd)
		if res == nil {
			log.Error("Close file number %d invalid", fd)
			return
		}

		f := res.Get().(*os.File)
		if f == nil {
			log.Warning("No open file for %d", fd)
		}else {
			err := f.Close()
			if err != nil {
				log.Error("Close file '%s' as #%d with error:%s", f.Name(), fd, err)
			}
		}
	case 50:// READ
		var err error
		fd := r1.Get()
		f := getFile(fd, fp)
		if f == nil {break}
		// TODO 是否可以直接使用 ReadAt 来避免 SEEK
		pos := r2.Get()
		_, err = f.Seek(int64(pos), os.SEEK_SET)
		if err != nil {goto READ_FAILED}
		switch o.Get(){
		case 16, 17:// float 和 int 操作相同
			b := make([]byte, 4)
			_, err = f.Read(b)
			if err != nil {goto READ_FAILED}
			r3.Set(int(int32(Codec.Uint32(b))))
		case 18:
			b := make([]byte, 1)
			str := make([]byte, 0)
			for {
				_, err = f.Read(b)
				if err != nil {goto READ_FAILED}
				if b[0] == 0 {break}
				str = append(str, b[0])
			}
			r3.SetStr(string(str))
		}
		break
		READ_FAILED:
		log.Error("Read failed '%s' as #%d at %d:%s", f.Name(), fd, pos, err.Error())
	case 51:// WRITE
		var err error
		fd := r1.Get()
		f := getFile(fd, fp)
		if f == nil {break}
		// TODO 是否可以直接使用 ReadAt 来避免 SEEK
		pos := r2.Get()
		_, err = f.Seek(int64(pos), os.SEEK_SET)
		if err != nil {goto WRITE_FAILED}
		switch o.Get(){
		case 16, 17:// float 和 int 操作相同
			b := make([]byte, 4)
			Codec.PutUint32(b, uint32(r3.Get()))
			_, err = f.Write(b)
			if err != nil {goto WRITE_FAILED}
		case 18:
			b := []byte(r3.Str())
			b = append(b, 0)
			_, err = f.Write(b)
			if err != nil {goto WRITE_FAILED}
		}
		break
		WRITE_FAILED:
		log.Error("Write failed '%s' as #%d at %d:%s", f.Name(), fd, pos, err.Error())
	case 52:
		fd := r3.Get()
		f := getFile(fd, fp)
		if f != nil {
			// Get file length
			ret, err := f.Seek(0, os.SEEK_CUR)
			if err != nil {
				log.Error("ftell '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}
			l, err := f.Seek(0, os.SEEK_END)
			if err != nil {
				log.Error("floc '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}
			_, err = f.Seek(ret, os.SEEK_SET)
			if err != nil {
				log.Error("fseek '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}

			r3.Set(0)
			if ret == l {
				r3.Set(1)
			}
		}
	case 53:
		fd := r3.Get()
		f := getFile(fd, fp)
		if f != nil {
			// Get file length
			ret, err := f.Seek(0, os.SEEK_CUR)
			if err != nil {
				log.Error("ftell '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}
			l, err := f.Seek(0, os.SEEK_END)
			if err != nil {
				log.Error("floc '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}
			_, err = f.Seek(ret, os.SEEK_SET)
			if err != nil {
				log.Error("fseek '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}
			r3.Set(int(l))
		}
	case 54:
		fd := r3.Get()
		f := getFile(fd, fp)
		if f != nil {
			ret, err := f.Seek(0, os.SEEK_CUR)
			if err != nil {
				log.Error("ftell '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}
			r3.Set(int(ret))
		}
	case 55:
		fd := r2.Get()
		f := getFile(fd, fp)
		if f != nil {
			_, err := f.Seek(int64(r3.Get()), 0)
			if err != nil {
				log.Error("Seek '%s' as %d faield: %s", f.Name(), fd, err.Error())
			}
		}
	}
}
func getFile(fd int, fp ResPool) (*os.File) {
	res := fp.Get(fd)
	if res == nil {
		log.Error("File number %d invalid", fd)
		return nil
	}
	i := res.Get()
	if i == nil {
		log.Warning("No open file for %d", fd)
	}
	f, ok := i.(*os.File)
	if !ok {
		log.Warning("No open file for %d", fd)
	}
	return f
}
func newFilePool() ResPool {
	p := &resPool{pool: make(map[int]Res), start:1, step:1, reuse:false, limit:10}
	// 预申请资源
	for i := 0; i < 10; i +=1 {
		p.Acquire()
	}
	return p
}
