$ErrorActionPreference = 'Stop'

$env:TRELLO_API_KEY = [Environment]::GetEnvironmentVariable('TRELLO_API_KEY', 'User')
$env:TRELLO_TOKEN = [Environment]::GetEnvironmentVariable('TRELLO_TOKEN', 'User')

if ([string]::IsNullOrWhiteSpace($env:TRELLO_API_KEY) -or [string]::IsNullOrWhiteSpace($env:TRELLO_TOKEN)) {
  Write-Error 'TRELLO_API_KEY e TRELLO_TOKEN precisam estar definidos nas variaveis de ambiente do usuario.'
  exit 1
}

& node 'C:\Users\Gustavo\Documents\Repositorios\IAs-MCPs-Skills\trello-cli\dist\index.js'
exit $LASTEXITCODE
