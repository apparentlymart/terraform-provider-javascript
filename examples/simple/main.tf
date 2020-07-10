
data "javascript" "example" {
  source = "num + num"
  vars = {
    num = 2
  }
}

output "source" {
  value = data.javascript.example.source
}

output "vars" {
  value = data.javascript.example.vars
}

output "result" {
  value = data.javascript.example.result
}
