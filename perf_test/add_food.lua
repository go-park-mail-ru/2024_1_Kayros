wrk.method = "PUT"
wrk.headers["Content-Type"] = "application/json"
wrk.headers["XCSRF-Token"] = "44f96d469525c8f5f62a664ec1b04514ec427ba407eb3a697e9c0efced84cd57.dcfde8c5-2041-4724-a033-50bde404ceab!S26JTfE6"
wrk.headers["Cookie"] = "session_id=dcfde8c5-2041-4724-a033-50bde404ceab; csrf_token=44f96d469525c8f5f62a664ec1b04514ec427ba407eb3a697e9c0efced84cd57.dcfde8c5-2041-4724-a033-50bde404ceab!S26JTfE6;"

local count = 1
local food_id = 72

request = function()
    count = count + 1

    local body = string.format('{"food_id": %d, "count": %d}', food_id, count)
    return wrk.format(nil, "/api/v1/order/food/update_count", nil, body)
end