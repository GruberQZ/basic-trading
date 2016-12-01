import json
import re
from tkinter import *
import requests

import tkinter as tk
import pprint
import urllib.request
import http.client

url = "https://b5b1f30cd80c4041972890286eb7e5df-vp0.us.blockchain.ibm.com:5001/chaincode"
#url = "http://maps.googleapis.com/maps/api/geocode/json?address=google"

 # data1 = str({"jsonrpc": "2.0","method": "query","params": {"type": 1,"chaincodeID":
 #     {"name": "76c1e5a0f389b61ed57ffb68be07aae7fa1c63dc98361afc50efb2d6fab41e1536c0270ffa31fb9a6b83d9829c33924d5f47ec16e3517df5c7dba80c082f758a"},"ctorMsg":
 #     {"function": "read","args": ["ece"]},"secureContext": "user_type1_1"},"id": 0})

#   method, params.ctorMsg.function, and params.ctorMsg.args are the only ones that can change

#print(data1)


data1 = "{\"jsonrpc\": \"2.0\",\n\"method\": " #method type here
data2 = "\n\"params\": {\n\"type\": 1,\n\"chaincodeID\":{\n\"name\": \"76c1e5a0f389b61ed57ffb68be07aae7fa1c63dc98361afc50efb2d6fab41e1536c0270ffa31fb9a6b83d9829c33924d5f47ec16e3517df5c7dba80c082f758a\"\n},\n\"ctorMsg\":{\n\"function\": " #params.ctorMsg.function args
data3 = "\n\"args\": [" # params.ctorMsg.args here
data4 = "]\n},\n\"secureContext\": \"user_type1_1\"\n},\n\"id\": 0\n}"


class Application(Frame):
    methodSel = IntVar()
    def __init__(self, master):

        Frame.__init__(self, master)

        self.grid()
        #self.pack()
        #self.pack_propagate(0)
        self.create_widgets()

    def create_widgets(self):

        self.textArea2 = Text(self, height = 10, width = 50)
        self.textArea2.grid(row = 0, column = 0, columnspan = 2)
        self.textArea2.insert(END, "This is where we will receive text from Bluemix")
        self.textArea2.config(state = DISABLED, font =("Times New Roman", 35))

        # self.label1 = Label(self, height = 2, text = "Method: ")
        # self.label1.grid(row = 1, column = 0, sticky = W)

        self.label2 = Label(self,  height = 2, text = "Function: ")
        self.label2.grid(row =2, column = 0, sticky = W)

        self.label3 = Label(self,  height = 2, text = "Arguments: ")
        self.label3.grid(row = 3, column = 0, sticky = W)

        self.button1 = Button(self, command = self.buttonClick)
        self.button1.grid(row = 4, column = 0, sticky = W)
        self.button1.config(text = "Send Command", height = 1, width = 20)

        # self.methodText = Text(self, height = 1, width = 20)
        # self.methodText.grid(row = 1, column = 1)

        self.functionText = Text(self, height = 1, width = 20)
        self.functionText.grid(row = 2, column = 1)

        self.argumentsText = Text(self, height = 1, width = 20)
        self.argumentsText.grid(row = 3, column = 1)

        self.queryButton = Radiobutton(self, font =("Ubuntu", 35), text = "Query", variable =self.methodSel, value = 1)
        self.invokeButton = Radiobutton(self, font =("Ubuntu", 35), text = "Invoke", variable =self.methodSel, value = 2)

        self.queryButton.config(command = self.radioSelect1(self.methodSel), borderwidth = 10)
        self.invokeButton.config(command = self.radioSelect1(self.methodSel), borderwidth = 10)

        self.queryButton.grid(row = 1, column = 0, sticky = W)
        self.invokeButton.grid(row = 1, column = 1, sticky = W)

        # self.textArea1 = Text(self, height = 20, width = 50)
        # self.textArea1.grid(row = 0, column = 0)
        # self.textArea1.insert(END, "This is Json sent to Bluemix")
        # self.textArea1.config(text = "This is where text will be sent to Bluemix")

    def radioSelect1(self, selection):
        print(selection.get(Application))
        if(selection == 1):
            print("QUERYING")
            self.methodSel = 0
            self.invokeButton.deselect()
            self.queryButton.select()
        if(selection == 2):
            print("INVOKING")
            self.methodSel = 1
            self.queryButton.deselect()
            self.invokeButton.select()


    def buttonClick(self):

        # method, params.ctorMsg.function, and params.ctorMsg.args are the only ones that can change

        # method = self.methodText.get("1.0", 'end-1c')
        function = self.functionText.get("1.0", 'end-1c')
        arguments = self.argumentsText.get("1.0", 'end-1c')

        print("MethodSel is: " + str(self.methodSel))

        if self.methodSel == 0:
            method = "query"
        else:
            method = "invoke"

        method = str("\"" + method + "\",")
        function = "\"" + function + "\","
        arguments = "\"" + arguments + "\""
        arguments = arguments.replace(',', '",\n"')

        data = data1 + method + data2 + function + data3 + arguments + data4
        # data = data.replace('\r\n', '\\r\\n')
        # data = data.replace('\t', '')
        print(data)
        try:
            r = requests.post(url, data = data, json = True)
            r.headers = 'Content-type', 'application/json'
            r.encoding = 'utf-8'
        except requests.ConnectionError as e:
            print(data)
            print("We messed up boys")
            print("Time to hang it up, let Kostas know we weren't read :'(")
            exit(e)

        response = r.text


        print(response)
        if "OK" not in response:
            print("We got a bad response")
            self.textArea2.config(state=NORMAL)
            self.textArea2.delete("1.0", END)
            self.textArea2.insert(END, "Entered Params are Invalid")
            self.textArea2.config(state=DISABLED)
            return


        response = response.replace('\\"', "")
        response = response.replace('"', "")
        response = response.replace('{status:OK,', "")
        if "jsonrpc:2.0," in response:
            response = response.replace("jsonrpc:2.0,", "")
        if ",id:0" in response:
            response = response.replace(",id:0", "")

        response = response.replace('{', "")
        response = response.replace('}', "")
        response = response.replace('message:', '\n\n')

        self.textArea2.config(state = NORMAL)
        self.textArea2.delete("1.0", END)
        self.textArea2.insert(END, response)
        self.textArea2.config(state = DISABLED)

        print("Data entered: \n\n" + data)
        print("Response received: \n\t" + r.text)


root = Tk()
root.title("Bluemix Blockchain-Interface")
root.geometry("1500x800")
app = Application(root)
root.mainloop()