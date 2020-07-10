# Terraform JavaScript Provider

This is a local-compute-only provider that allows to to evaluate JavaScript
code via a Terraform data source.

The JavaScript execution environment is minimal and does not allow importing
external libraries, etc. Instead, this is intended for situations where you
want to process a complex data structure and it ends up being more concise
to express the algorithm as a small JavaScript program rather than as
one or more Terraform language expressions.

This provider supports ECMAScript 5.1 and, due to the intended use for
wrangling data structures, [Underscore.js](https://underscorejs.org/)
1.10.2 is automatically available to enable concise programming in functional
style.

When considering use of this provider, keep the following advice in mind:
_just because you can, it doesn't mean you should_. If something already
has a relatively straightforward and intuitive representation in the Terraform
language, you can make your configuration simpler and more accessible by
staying within the Terraform language. Use this provider only for the rare
problems were some complex computation is required and JavaScript features
such as higher-order functions would make the result more concise and easier
to read and maintain.

This is not a HashiCorp-maintained provider.

## Usage

To make the provider available for use in your module, you must first declare
it as a required provider inside your `terraform` block:

```hcl
terraform {
  required_providers {
    javascript = {
      source = "apparentlymart/javascript"

      # Until this provider has a stable release, always select an exact
      # version to avoid adopting new releases that may have breaking changes.
      version = "0.0.1"
    }
  }
}
```

Once the provider is available in your module you can write one or more
`javascript` data resources:

```hcl
data "javascript" "example" {
  source = "num + num"
  vars = {
    num = 2
  }
}

output "result" {
  value = data.javascript.example.result
}
```

The `source` argument specifies JavaScript source code to compile and execute.
If you have a longer program you can either use Terraform's flush heredoc
syntax or you can load the program from an external file using Terraform's
built-in function
[`file`](https://www.terraform.io/docs/configuration/functions/file.html).

You can optionally set `vars` to a mapping, in which case each element of
the mapping will become a global variable in the JavaScript scope.

The data source exports an attribute called `result`, which is the result of
the final expression evaluated by the program.

## Input and Output Value Mapping

Because Terraform and JavaScript have separate type systems, this provider
must therefore translate values from the `vars` mapping into JavaScript and
translate the result value from JavaScript.

The mapping for input values in `vars` is effectively the same as for
Terraform's built in
[`jsonencode`](https://www.terraform.io/docs/configuration/functions/jsonencode.html)
function, as if the result of that function had been then been parsed as JSON
within the JavaScript context.

JavaScript's type system is larger than Terraform's and includes concepts that
don't exist in Terraform, such as function values. In order to achieve a
well-defined and easy-to-understand mapping from JavaScript to Terraform, this
provider interprets the result the same way as JavaScript's own `JSON.stringify`
function would, including calling a `toJSON` method when available, and then
applies the same mapping as for Terraform's built in
[`jsondecode`](https://www.terraform.io/docs/configuration/functions/jsondecode.html)
function.

The `JSON.stringify` function reacts to some unsupported values by replacing
them with `null` and to others by returning an error. This provider's
interpretation of the result follows those same conventions, and will raise
an error in the same situations that `JSON.stringify` would.

## Requirements

This provider requires at least Terraform 0.12 and can only be installed
automatically by Terraform 0.13 or later.
