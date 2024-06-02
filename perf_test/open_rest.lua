wrk.method = "GET"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["XCSRF-Token"] = "380869aba9e630b7625b3226164f387552185e9e69f9fdff62ae008febcf9d23.ff361a24-621c-4854-b416-fcbe6b1122f2!ZE8T8yEW"
wrk.headers["Cookie"] = "session_id=b09ccf30-d099-480f-aa8c-71a04a049d44; csrf_token=c13ca394c9e3f18b1c75ec35e6fa8dd90db3aeac6f96fd08b290d5e3001ae58e.b09ccf30-d099-480f-aa8c-71a04a049d44!CYhR5ej0;"

request = function()
    return wrk.format(nil, "/api/v1/restaurants/4", nil, nil)
end