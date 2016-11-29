import json
from tkinter import *
import pprint
import urllib.request
import http.client

url = "https://b5b1f30cd80c4041972890286eb7e5df-vp0.us.blockchain.ibm.com:5001/chaincode"
data1 = str({"jsonrpc": "2.0","method": "query","params": {"type": 1,"chaincodeID":
    {"name": "76c1e5a0f389b61ed57ffb68be07aae7fa1c63dc98361afc50efb2d6fab41e1536c0270ffa31fb9a6b83d9829c33924d5f47ec16e3517df5c7dba80c082f758a"},"ctorMsg":
    {"function": "read","args": ["ece"]},"secureContext": "user_type1_1"},"id": 0})

print(data1)

class Application(Frame):

    def __init__(self, master):

        Frame.__init__(self, master)
        self.grid()
        self.create_widgets()

    def create_widgets(self):

        self.button1 = Button(self, text = "Hit this button to send stuff to Bluemix", command=self.buttonClick)
        self.button1.grid(row = 1, column = 0, sticky = S)
        self.button1.config(text = "Click here to send json to Bluemix", height = 1, width = 40)


        self.textArea1 = Text(self, height = 20, width = 50)
        self.textArea1.grid(row = 0, column = 0)
        self.textArea1.insert(END, "This is Json sent to Bluemix")
        #self.textArea1.config(text = "This is where text will be sent to Bluemix")

        self.textArea2 = Text(self, height = 20, width = 50)
        self.textArea2.grid(row = 0, column = 1)
        self.textArea2.insert(END, "This is where we will receive text from Bluemix")
        self.textArea2.config(state = DISABLED)

    def buttonClick(self):

    #   method, params.ctorMsg.function, and params.ctorMsg.args are the only ones that can change

        data = self.textArea1.get("1.0", END)
        data = json.loads(data)
        req = urllib.request.Request(url)
        req.add_header('Content-type', 'application/json')
        response = urllib.request.urlopen(req, json.dumps(data).encode('utf-8'))
        httpResponse =http.client.HTTPResponse.read(response)

        self.textArea2.config(state = NORMAL)
        self.textArea2.delete("1.0", END)
        self.textArea2.insert(END, httpResponse)
        self.textArea2.config(state = DISABLED)

        pprint.pprint(httpResponse)

root = Tk()
root.title("Send json to Bluemix")
root.geometry("800x400")
app = Application(root)
root.mainloop()