if not os.stat('/etc/sv/nsqlookupd') then return end

local nsqlookupd = { host = "http://127.0.0.1:4161" }

function nsqlookupd.get(url)
  print(nsqlookupd)
  print(nsqlookupd.host)
  local result, err = http.get(nsqlookupd.host..url)
  if err then metrics.set("nsqlookupd.problem", url.." availability: "..err); return nil; end
  if not (result.code == 200) then metrics.set("nsqlookupd.problem", url.." response code: "..tostring(result.code)); return nil; end
  return result.body
end

function nsqlookupd.ping()
  local body = nsqlookupd.get("/ping")
  if body == nil then return end
  if not (body == "OK") then
    metrics.set("nsqlookupd.problem", "/ping response body: "..tostring(body))
    return false
  end
  metrics.set("nsqlookupd.problem", "OK")
  return true
end

function nsqlookupd.check_nodes()
  local body = nsqlookupd.get("/nodes")
  if body == nil then return end
  local response = json.decode(body)

  local producers = 0
  for _, producer in pairs(response.producers) do
    producers = producers + 1
  end
  metrics.set("nsqlookupd.nodes.producers", producers)
end

-- бесконечный loop
while true do
  if nsqlookupd.ping() then
    nsqlookupd.check_nodes()
  end
  time.sleep(60)
end
