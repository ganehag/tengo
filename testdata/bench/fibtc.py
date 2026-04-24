import sys
sys.setrecursionlimit(100000)

def fib(x, a=0, b=1):
    if x == 0: return a
    if x == 1: return b
    return fib(x-1, b, a+b)

print(fib(35))
