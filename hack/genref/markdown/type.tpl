{{ define "type" }}

## `{{ .Name.Name }}`     {#{{ .Anchor }}}
    
{{ if eq .Kind "Alias" -}}
(Alias of `{{ .Underlying }}`)
{{ end }}

{{- with .References }}
**Appears in:**
{{ range . }}
{{ if or .Referenced .IsExported -}}
- [{{ .DisplayName }}]({{ .Link }})
{{ end -}}
{{- end -}}
{{- end }}

{{ if .GetComment -}}
{{ .GetComment }}
{{ end }}
{{ if .GetMembers -}}

<table class="table">
<thead><tr><th width="30%">Field</th><th>Description</th></tr></thead>
<tbody>
    {{- if .IsExported -}}
<tr><td><code>apiVersion</code><br/>string</td><td><code>{{- .APIGroup -}}</code></td></tr>
<tr><td><code>kind</code><br/>string</td><td><code>{{- .Name.Name -}}</code></td></tr>
    {{ end -}}

{{/* The actual list of members is in the following template */}}
{{- template "members" . -}}
</tbody>
</table>
{{- end -}}
{{- end -}} 