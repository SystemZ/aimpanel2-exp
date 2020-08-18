package host

import (
	"bytes"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/ecode"
	"gitlab.com/systemz/aimpanel2/master/config"
	"gitlab.com/systemz/aimpanel2/master/model"
	"text/template"
)

type DeploymentOptions struct {
	Token   string
	UrlRepo string
	UrlApi  string
}

func GenerateDeploymentScript(hostToken string) (err error, res string) {
	// FIXME throw template error if some variable is missing
	//tmpl := template.New("main").Option("missingkey=error")
	//_, err := tmpl.Parse(templateToParse)

	host, err := model.GetHostByToken(hostToken)
	if err != nil {
		err = &lib.Error{ErrorCode: ecode.HostNotFound}
		return
	}
	httpParams := DeploymentOptions{
		Token:   host.Token,
		UrlRepo: config.HTTP_REPO_URL,
		UrlApi:  config.HTTP_API_URL,
	}
	t, err := template.ParseFiles(config.HTTP_TEMPLATE_DIR + "install.sh")
	if err != nil {
		err = &lib.Error{ErrorCode: ecode.TemplateGenerationFailed}
		return
	}
	var templateOutput bytes.Buffer
	err = t.Execute(&templateOutput, httpParams)
	if err != nil {
		// FIXME other error code for template exec
		err = &lib.Error{ErrorCode: ecode.TemplateGenerationFailed}
		return
	}
	res = templateOutput.String()
	return
}
