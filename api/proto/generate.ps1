$protos = Get-ChildItem -Path "evidra/v1/*.proto" -Name
$includes = ".", "$env:TEMP\protoc\include"

protoc `
  --proto_path="." `
  --go_out="../../api/gen" `
  --go_opt="paths=source_relative" `
  --go-grpc_out="../../api/gen" `
  --go-grpc_opt="paths=source_relative" `
  $protos

Write-Output "Generated gRPC code for: $protos"
