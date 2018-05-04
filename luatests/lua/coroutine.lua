local function cof(x)
    print("in cof", x)
    print("in cof", coroutine.yield(x + 2))
    return "from cof"
end

local co = coroutine.create(cof)
print("out", coroutine.resume(co, 1))
print("out", coroutine.resume(co, "two"))

--> =in cof	1
--> =out	true	3
--> =in cof	two
--> =out	true	from cof

print(coroutine.running())
--> ~^thread:.*\ttrue$

do
    local function cof()
        return coroutine.running()
    end
    local co = coroutine.create(cof)

    print(coroutine.resume(co))
--> ~^true\tthread:.*\tfalse$
end
