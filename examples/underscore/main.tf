
data "javascript" "example" {
  source = <<-EOT
    _.reduce(input, function(memo, num) { return memo + num; }, 0)
  EOT
  vars = {
    input = [1, 2, 3]
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
