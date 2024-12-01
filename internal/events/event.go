package events

type Event string

const (
	ShowLoginForm      Event = "showLoginForm"
	ShowMainMenu       Event = "showMainMenu"
	ShowRegisterForm   Event = "showRegisterForm"
	ShowAddDataMenu    Event = "showAddDataMenu"
	ShowAddDataForm    Event = "showAddDataForm"
	ShowUpdateDataForm Event = "showUpdateDataForm"
	ShowEditDataForm   Event = "showEditDataForm"
	ShowGetData        Event = "showGetData"
	ShowDataDetails    Event = "showDataDetails"
	ShowRemoveDataForm Event = "showRemoveDataForm"
)

type AddDataFormEvent struct {
	Type string
}

type EditDataFormEvent struct {
	ID         uint64
	DataStruct map[string]interface{}
}

type DataDetailsFormEvent struct {
	ID          uint64
	DataType    string
	Metadata    string
	DisplayData string
}
