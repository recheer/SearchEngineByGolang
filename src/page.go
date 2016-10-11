package page

import (
  "fmt"
  "unsafe"
  "container/list"
)

const maxkeylen uint32 = 64
const pagesize int64 = 1024 * 4 * 2

const (
	tindex     uint8 = 0
	tdata      uint8 = 1
)

type dictitem struct {
  word [maxkeylen]byte
  wordid uint32
}

type indexitem struct {
  wordid uint32
  docid uint32
}

type index struct {
  wordid uint32
  count uint32
  offset uint32     //the offset of key in data file offset+sizeof(uint32)
}

type dataitem uint32

type page struct {
  pgid uint32
  nextid uint32
  preid uint32
  pgtype uint8       //dacid data page and index page
  curSize uint32
  data uintptr      //save data for docid or index page
}

func (p *page) getDataPtr {
  return p.ptype ? (*uint32)(unsafe.Pointer(&p.data)) : &((*[0xFFFF]item)(unsafe.Pointer(&p.data)))[0]
}

func (p *page) insertItem(i indexkv) bool {
  if p.curSize + i.v.len() + maxkeylen > pagesize {
    return false
  }

}

func (bt *BTreedb) Newpage( /*parent, pre *page*/ parentid, preid uint32, pagetype uint8) (*page, *page, *page) {

	if bt.checkmmap() != nil {
		fmt.Printf("check error \n")
		return nil, nil, nil
	}
	var parent *page
	var pre *page
	lpage := (*page)(unsafe.Pointer(&bt.mmapbytes[(int64(bt.meta.maxpgid) * pagesize)]))
	//fmt.Printf("lapge:%v\n", unsafe.Pointer(lpage))
	lpage.curid = bt.meta.maxpgid
	lpage.pgtype = pagetype
	lpage.nextid = 0
	lpage.preid = 0
	if pagetype == tinterior {
		lpage.count = 1
		ele := (*[0xFFFF]element)(unsafe.Pointer(&bt.mmapbytes[(int64(bt.meta.maxpgid)*pagesize + pageheadOffset)]))
		lpage.used = uint32(pageheaadlen)
		ele[0].setkv("", 0)
		lpage.elementsptr = uintptr(unsafe.Pointer(ele))

	} else {
		lpage.count = 0
		ele := (*[0xFFFF]element)(unsafe.Pointer(&bt.mmapbytes[(int64(bt.meta.maxpgid)*pagesize + pageheadOffset)]))
		lpage.elementsptr = uintptr(unsafe.Pointer(ele))
		lpage.used = uint32(pageheaadlen)
	}
	//fmt.Printf("lapge:%v\n", unsafe.Pointer(lpage))
	//fmt.Printf("parent:%v\n", unsafe.Pointer(parent))
	if parentid != 0 {
		parent = bt.getpage(parentid)
		lpage.parentpg = parent.curid
	} else {
		lpage.parentpg = 0
	}

	if preid != 0 {
		pre = bt.getpage(preid)
		lpage.nextid = pre.nextid
		pre.nextid = lpage.curid
		lpage.preid = pre.curid
	}

	bt.meta.maxpgid++
	return lpage, parent, pre
}

func (bt *BTreedb) checkmmap() error {
	if int(int64(bt.meta.maxpgid)*pagesize) >= len(bt.mmapbytes) {
		err := syscall.Ftruncate(int(bt.fd.Fd()), int64(bt.meta.maxpgid+1)*pagesize)
		if err != nil {
			fmt.Printf("ftruncate error : %v\n", err)
			return err
		}
		maxpgid := bt.meta.maxpgid
		syscall.Munmap(bt.mmapbytes)
		//fmt.Printf(".meta.maxpgid:%v\n",bt.meta.maxpgid)
		bt.mmapbytes, err = syscall.Mmap(int(bt.fd.Fd()), 0, int(int64( /*bt.meta.maxpgid*/ maxpgid+1)*pagesize), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED)

		if err != nil {
			fmt.Printf("MAPPING ERROR  %v \n", err)
			return err
		}

		bt.meta = (*metaInfo)(unsafe.Pointer(&bt.mmapbytes[0]))

	}
	return nil
}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) (bool) {
 var exist = true;
 if _, err := os.Stat(filename); os.IsNotExist(err) {
  exist = false;
 }
 return exist;
}

func writeDictFile(filename string, key []byte, value uint32) {
  if checkFileIsExist(filename) {  //如果文件存在
    f, _ = os.OpenFile(filename, os.O_APPEND, 0666)  //打开文件
    fmt.Println("file exist!");
   }else {
    f, _ = os.Create(filename)  //创建文件
    fmt.Println("file is not exist!");
   }
   /*strkey := []byte(key)
   if len(key) > maxkeylen {
     strkey = strings(key[0:maxkeylen-1])
   }*/
   strkey = strings(key)
   strval := strconv.Itoa(value)
   str := strings.Join({strkey, strval}, ",")
   err := ioutil.WriteFile(filename, str, 0666)  //写入文件(字节数组)
   if !err {
     fmt.Println("write file error!");
   }
}

func loadDictFile(filename string) map[string]uint32 {
  dict := make(map[string]uint32)
}
