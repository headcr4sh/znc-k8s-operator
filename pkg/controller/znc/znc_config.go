package znc

import (
	"bytes"
	"text/template"
	zncv1 "znc-operator/pkg/apis/znc/v1"
)

func configurationTemplate() string {
	return `// Automatically generated configuration file.
Version = {{ .GetVersion }}

<Listener irc>
        AllowIRC = true
        AllowWeb = false
        Port = 6667
        IPv4 = true
        IPv6 = false
        SSL = false
</Listener>

<Listener web>
        AllowIRC = false
        AllowWeb = true
        Port = 8080
        IPv4 = true
        IPv6 = false
        SSL = false
</Listener>

{{- with .GetConfig }}
{{- range .LoadModules }}
LoadModule = {{ . }}
{{- end }}

AnonIPLimit = {{ .AnonIPLimit }}
ConnectDelay = {{ .ConnectDelay }}
HideVersion = {{ .HideVersion }}
MaxBufferSize = {{ .MaxBufferSize }}

{{- range .Motd }}
Motd = {{ . }}
{{- end }}

ServerThrottle = {{ .ServerThrottle }}
StatusPrefix = {{ .GetStatusPrefix }}

{{- range .Users }}
<User {{ .Name }}>
        Admin = {{ .Admin }}
        Allow = *
        AltNick = {{ .AltNick }}
        AppendTimestamp = {{ .AppendTimestamp }}
        AutoClearChanBuffer = {{ .AutoClearChanBuffer }}
        AutoClearQueryBuffer = {{ .AutoClearQueryBuffer }}
        Buffer = {{ .Buffer }}
        ChanBufferSize = {{ .ChanBufferSize }}
        ChanModes = {{ .GetChanModes }}
        ClientEncoding = {{ .GetClientEncoding }}
        Ident = {{ .GetIdent }}
        JoinTries = {{ .JoinTries }}
        {{- range .LoadModules }}
        LoadModule = {{ . }}
        {{- end }}
        MaxJoins = {{ .MaxJoins }}
        MaxQueryBuffers = {{ .MaxQueryBuffers }}
        MultiClients = {{ .MultiClients }}
        Nick = {{ .Nick }}
        NoTrafficTimeout = {{ .NoTrafficTimeout }}
        PrependTimestamp = {{ .PrependTimestamp }}
        QueryBufferSize = {{ .QueryBufferSize }}
        {{- if .QuitMsg }}QuitMsg = {{ .QuitMsg }}{{ end }}
        {{- if .RealName }}RealName = {{ .RealName }}{{ end }}
        StatusPrefix = {{ .GetStatusPrefix }}
        {{- if .TimestampFormat }}TimestampFormat = {{ .TimestampFormat }}{{ end }}
        {{- if .Timezone }}Timezone = {{ .Timezone }}{{ end }}
        Pass = {{ .Pass }}
        {{- range .Networks }}
        <Network {{ .Name }}>
                {{- if .AltNick }}AltNick = {{ .AltNick }}{{ end }}
                {{- if .Encoding }}Encoding = {{ .Encoding }}{{ end }}
                {{- if .Ident }}Ident = {{ .Ident }}{{ end }}
                IRCConnectEnabled = {{ .IRCConnectEnabled }}
                {{- if .JoinDelay }}JoinDelay = {{ .JoinDelay }}{{ end }}
                {{- range .LoadModules }}
                LoadModule = {{ . }}
                {{- end }}
                {{- if .Nick }}Nick = {{ .Nick }}{{ end }}
                {{- if .QuitMsg }}QuitMsg = {{ .QuitMsg }}{{ end }}
                {{- if .RealName }}RealName = {{ .RealName }}{{ end }}
                {{- range .Servers }}
                Server = {{ . }}
                {{- end }}
                {{- range .Channels }}
                <Channel {{ .Name }}>
                        AutoClearChanBuffer = {{ .AutoClearChanBuffer }}
                        Buffer = {{ .Buffer }}
                        Detached = {{ .Detached }}
                        Disabled = {{ .Disabled }}
                        {{- if .Key }}Key = {{ .Key }}{{ end }}
                        {{- if .Modes }}Modes = {{ .Modes }}{{ end }}
                </Channel>
                {{- end }}
        </Network>
        {{- end }}
</User>
{{- end }}

{{- end }}
`
}

func RenderConfiguration(spec *zncv1.ZNCSpec) (cfg string, err error) {
	tmpl, err := template.New("znc-config").Parse(configurationTemplate())
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, spec)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
