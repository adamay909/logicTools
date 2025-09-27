#!/bin/lua


local indexSkel = "html/indexSkel.html"
local body = "html/body.html"
local readme = "html/README.html"
local indexFile = "html/index.html"
local help = "html/help.html"

-- Create a dummy file for demonstration

-- 1. Open for reading
local file_read = io.open(readme, "r")
if not file_read then
    error("Failed to open file for reading.")
end

-- 2. Read content
local readmestr = file_read:read("*all")
file_read:close()

local file_read = io.open(help, "r")
if not file_read then
    error("Failed to open file for reading.")
end

-- 2. Read content
local helpstr = file_read:read("*all")
file_read:close()

local file_read = io.open(body, "r")
if not file_read then
    error("Failed to open file for reading.")
end

-- 2. Read content
local bodystr = file_read:read("*all")
file_read:close()


-- 3. Modify content
local bodystr = bodystr:gsub("#README#", readmestr)
local bodystr = bodystr:gsub("#HELP#", helpstr)

local file_read = io.open(indexSkel, "r")
if not file_read then
    error("Failed to open file for reading.")
end

-- 2. Read content
local indexstr = file_read:read("*all")
file_read:close()

local indexstr = string.gsub(indexstr, "#BODY#", bodystr)

-- 4. Open for writing (overwrites existing file)
local file_write = io.open(indexFile, "w")
if not file_write then
    error("Failed to open file for writing.")
end

-- 5. Write modified content
file_write:write(indexstr)
file_write:close()

