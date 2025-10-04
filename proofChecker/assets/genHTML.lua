#!/bin/lua


local indexSkel = "html/indexSkel.html"
local body = "html/body.html"
local readme = "html/README.html"
local indexFile = "html/index.html"
local help = "html/help.html"
local branch = "html/branch"


local file_read = io.open(readme, "r")
if not file_read then
    error("Failed to open file for reading.")
end
local readmestr = file_read:read("*all")
file_read:close()

local file_read = io.open(help, "r")
if not file_read then
    error("Failed to open file for reading.")
end
local helpstr = file_read:read("*all")
file_read:close()

local file_read = io.open(body, "r")
if not file_read then
    error("Failed to open file for reading.")
end
local bodystr = file_read:read("*all")
file_read:close()

local bodystr = bodystr:gsub("#README#", readmestr)
local bodystr = bodystr:gsub("#HELP#", helpstr)

local file_read = io.open(indexSkel, "r")
if not file_read then
    error("Failed to open file for reading.")
end

local indexstr = file_read:read("*all")
file_read:close()

local indexstr = string.gsub(indexstr, "#BODY#", bodystr)

local file_read = io.open(branch, "r")
if file_read then 
 local branchname = file_read:read("*all")
 file_read:close()
indexstr = string.gsub(indexstr, "<title>Proof Checker</title>", "<title> Proof Checker ("..branchname..")</title>")
end

local file_write = io.open(indexFile, "w")
if not file_write then
    error("Failed to open file for writing.")
end

file_write:write(indexstr)
file_write:close()

