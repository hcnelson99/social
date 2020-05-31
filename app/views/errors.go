package views

func GetError(view *viewState) {
	view.Templates.ExecuteTemplate(view.response, "error.tmpl", nil)
}
