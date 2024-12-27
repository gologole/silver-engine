local http = require("socket.http")
local ltn12 = require("ltn12")
local json = require("dkjson")
local luaunit = require("luaunit")

local BASE_URL = "http://localhost:8080/api"

local function request(method, path, body)
    local response_body = {}
    local request_body = body and json.encode(body) or ""
    local headers = {
        ["Content-Type"] = "application/json",
        ["Content-Length"] = #request_body
    }
    local response, code, response_headers = http.request {
        url = BASE_URL .. path,
        method = method,
        headers = headers,
        source = ltn12.source.string(request_body),
        sink = ltn12.sink.table(response_body)
    }
    return code, json.decode(table.concat(response_body))
end

local function assert_status(expected, actual)
    luaunit.assertEquals(actual, expected, string.format("Expected status %d, got %d", expected, actual))
end

TestBooks = {}

function TestBooks:test_get_books()
    local code, body = request("GET", "/books")
    assert_status(200, code)
    luaunit.assertIsTable(body)
end

function TestBooks:test_add_book()
    local new_book = {title = "Test Book", author = "Test Author"}
    local code, body = request("POST", "/books", new_book)
    assert_status(201, code)
    luaunit.assertEquals(body.title, new_book.title)
    luaunit.assertEquals(body.author, new_book.author)
    return body.id
end

function TestBooks:test_delete_book()
    local book_id = self:test_add_book()
    local code, body = request("DELETE", "/books/" .. book_id)
    assert_status(200, code)
    luaunit.assertEquals(body.message, "Book deleted")
end

TestUsers = {}

function TestUsers:test_get_users()
    local code, body = request("GET", "/users")
    assert_status(200, code)
    luaunit.assertIsTable(body)
end

function TestUsers:test_add_user()
    local new_user = {name = "Test User", email = "test@example.com"}
    local code, body = request("POST", "/users", new_user)
    assert_status(201, code)
    luaunit.assertEquals(body.name, new_user.name)
    luaunit.assertEquals(body.email, new_user.email)
    return body.id
end

function TestUsers:test_delete_user()
    local user_id = self:test_add_user()
    local code, body = request("DELETE", "/users/" .. user_id)
    assert_status(200, code)
    luaunit.assertEquals(body.message, "User deleted")
end

TestLoans = {}

function TestLoans:test_get_loans()
    local code, body = request("GET", "/loans")
    assert_status(200, code)
    luaunit.assertIsTable(body)
end

function TestLoans:test_issue_loan()
    local user_id = TestUsers:test_add_user()
    local book_id = TestBooks:test_add_book()
    local new_loan = {user_id = user_id, book_id = book_id}
    local code, body = request("POST", "/loans", new_loan)
    assert_status(201, code)
    luaunit.assertEquals(body.user_id, new_loan.user_id)
    luaunit.assertEquals(body.book_id, new_loan.book_id)
    return body.id
end

function TestLoans:test_return_loan()
    local loan_id = self:test_issue_loan()
    local code, body = request("POST", "/loans/" .. loan_id .. "/return")
    assert_status(200, code)
    luaunit.assertNotNil(body.return_date)
end

TestStatistics = {}

function TestStatistics:test_get_statistics()
    local code, body = request("GET", "/statistics")
    assert_status(200, code)
    luaunit.assertIsTable(body)
end

os.exit(luaunit.LuaUnit.run())