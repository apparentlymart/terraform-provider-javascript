
data "javascript" "example" {
  source = file("${path.module}/fibonacci.js")
  vars = {
    input = 10
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
