x = 0
while true do
    metrics.set("x", x)
    time.sleep(1)
    x = x + 1
    if x > 3 then error("timeout in plugin") end
end
