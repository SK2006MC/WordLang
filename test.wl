from stdio import print
from stdio import println

let myNumber be 10 # Changed "my number" to "myNumber"
let greeting be "Hello WordLang!"

print greeting

print add myNumber and 5 # returns the total ,mynumber not modified
print add 5 to myNumber # myNumber += 5 returns the total

if myNumber greater than 5 then # Changed "my number" to "myNumber"
  println "Number is greater than 5"
else
  println "Number is not greater than 5"
endif

let count be 0
while count less than 3 do
  println "Count:" count
  #let count be add count and 1 or
  #increment count by 1 or
  increment count # defaults to 1
endwhile

let names be strings "Alice" "Bob" "Charlie"
foreach name in names do
  print "Hello" name
endforeach

function greet person name
  print "Greetings" person
  print "Your name is" name
endfunction

call greet "User" "WordLang Learner"

let myList be numbers 1 2 3 4  # Changed "mylist" to "myList"
#let firstItem be get item at index 0 from myList # 
let firstItem be item at index 0 myList # need to find a best way
print firstItem

let isTen be isdefined mynumber # mynumber is not defined, should be false
print isTen
