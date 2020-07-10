function fib(num) {
  var ret = [];

  for (var i = 0; i < num; i++) {
    var a = ret[i - 1] || 0;
    var b = ret[i - 2] || 1;
    ret.push(a + b);
  }

  return ret;
}

fib(input);
