import json
from tkinter import *
import tkinter as tk
import pprint
import urllib.request
import http.client

url = "https://b5b1f30cd80c4041972890286eb7e5df-vp0.us.blockchain.ibm.com:5001/chaincode"
url = "http://www.google.com"

 # data1 = str({"jsonrpc": "2.0","method": "query","params": {"type": 1,"chaincodeID":
 #     {"name": "76c1e5a0f389b61ed57ffb68be07aae7fa1c63dc98361afc50efb2d6fab41e1536c0270ffa31fb9a6b83d9829c33924d5f47ec16e3517df5c7dba80c082f758a"},"ctorMsg":
 #     {"function": "read","args": ["ece"]},"secureContext": "user_type1_1"},"id": 0})

#   method, params.ctorMsg.function, and params.ctorMsg.args are the only ones that can change

#print(data1)

data1 = "{\"jsonrpc\": \"2.0\",\"method\": " #method type here
data2 = ",\"params\": {\"type\": 1,\"chaincodeID\":{\"name\": \"76c1e5a0f389b61ed57ffb68be07aae7fa1c63dc98361afc50efb2d6fab41e1536c0270ffa31fb9a6b83d9829c33924d5f47ec16e3517df5c7dba80c082f758a\"},\"ctorMsg\":{\"function\": " #params.ctorMsg.function args
data3 = ",\"args\": [" # params.ctorMsg.args here
data4 = "]},\"secureContext\": \"user_type1_1\"},\"id\": 0}"



class Application(Frame):

    def __init__(self, master):

        Frame.__init__(self, master)
        #self.grid()
        self.pack()
        #self.pack_propagate(0)
        self.create_widgets()

    def create_widgets(self):

        self.textArea2 = Text(self, height = 15, width = 50)
        self.textArea2.grid(row = 0, column = 0, columnspan = 2)
        self.textArea2.insert(END, "This is where we will receive text from Bluemix")
        self.textArea2.config(state = DISABLED)

        self.label1 = Label(self, height = 2, text = "Method: ")
        self.label1.grid(row = 1, column = 0, sticky = W)

        self.label2 = Label(self,  height = 2, text = "Function: ")
        self.label2.grid(row =2, column = 0, sticky = W)

        self.label3 = Label(self,  height = 2, text = "Arguments: ")
        self.label3.grid(row = 3, column = 0, sticky = W)

        self.button1 = Button(self, command = self.buttonClick)
        self.button1.grid(row = 4, column = 0, sticky = W)
        self.button1.config(text = "Send json to Bluemix", height = 1, width = 20)

        self.methodText = Text(self, height = 1, width = 20)
        self.methodText.grid(row = 1, column = 1)

        self.functionText = Text(self, height = 1, width = 20)
        self.functionText.grid(row = 2, column = 1)

        self.argumentsText = Text(self, height = 1, width = 20)
        self.argumentsText.grid(row = 3, column = 1)

        # self.textArea1 = Text(self, height = 20, width = 50)
        # self.textArea1.grid(row = 0, column = 0)
        # self.textArea1.insert(END, "This is Json sent to Bluemix")
        # self.textArea1.config(text = "This is where text will be sent to Bluemix")

    def buttonClick(self):

    #   method, params.ctorMsg.function, and params.ctorMsg.args are the only ones that can change

        method = self.methodText.get("1.0", END)
        function = self.functionText.get("1.0", END)
        arguments = self.argumentsText.get("1.0", END)

        method = "\"" + str(method) + "\""
        function = "\"" + str(function) + "\""
        arguments = "\"" + str(arguments) + "\""

        #print(data1 + method + data2)

        data = data1 + method + data2 + function + data3 + arguments + data4
        print(data)
        #data = data.replace('\r\n', '\\r\\n')

        #print("json entered was: " + data)

        #print(data)

        data = json.loads(data, strict = False)

        #print(data)

        req = urllib.request.Request(url)
        req.add_header('Content-type', 'application/json')

        # Test opening of google
        try:
            urllib.request.urlopen(url)
        except urllib.request.HTTPError as e:
            print(e.code)
            print(e.msg)
            print("We couldn't get to: " + str(url))
            print("Time to kill yourself...")
            exit()
        print("Back from opening url")

        response = urllib.request.urlopen(req, json.dumps(data).encode('utf-8'))

        httpResponse =http.client.HTTPResponse.read(response)

        self.textArea2.config(state = NORMAL)
        self.textArea2.delete("1.0", END)
        self.textArea2.insert(END, httpResponse)
        self.textArea2.config(state = DISABLED)

        pprint.pprint(httpResponse)

root = Tk()
root.title("Bluemix Blockchain-Interface")
root.geometry("450x425")
app = Application(root)
root.mainloop()