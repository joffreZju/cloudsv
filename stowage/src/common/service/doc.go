package service

//隐藏文档，启用文档
func SetDocStatus(st, tp, did int) (err error) {
	if st == model.DocUsing {
		err = model.SetDocHide(tp)
		if err != nil {
			return err
		}
	}

	d := &model.Document{
		Status: st,
		Id:     did,
	}
	err = model.UpdateDocumnet(d, "Status")
	return
}

func NewDocument(d *Document) (err error) {
	model.SetDocHide(d.DocType)
	err = model.CreateDocumnet(d)
	return
}

func NewFile(f *File) (err error) {
	err = model.CreateFile(f)
	return
}
