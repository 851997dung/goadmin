#!/usr/bin/python3

# Open a file
import string
 
from asyncio.windows_events import NULL


fo = open("./flowtype.dsl", "r",encoding='UTF-8')
print ("Name of the file: ", fo.name)
temp = fo.readlines()

res2 = "package flowDef\n"
res2+="var FlowTypeString map[int]string\n"
res2+="const (\n"
res1 ="func GetFlowTypeStrings()map[int]string{\n"
res1 += "FlowTypeString = make(map[int]string)\n"



for s in temp:        
   res2 += s
   s = s.rstrip("\n")
   temp1 = s.split("//")

   if temp1[0] == "LFT_NONE=iota":
       continue 
   if len(temp1) < 3 :
    if temp1[0]!= None:
        temp1[0] = temp1[0].rstrip(string.digits)
        temp1[0]= temp1[0].rstrip('=')
        res1 += "FlowTypeString["+temp1[0]+"] = "
    if temp1[1]!=None:
       res1+="gotext.Get(\""+temp1[1] +"\")"+"\n"


res2 += "\n)\n"
res1 +="return FlowTypeString \n}"

fo1 = open("../../common/def/flowDef/logFlowType.go","w",encoding='UTF-8')

fo1.write(res2)
fo1.write(res1)