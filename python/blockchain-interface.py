from tkinter import *
import requests

url = "https://b5b1f30cd80c4041972890286eb7e5df-vp0.us.blockchain.ibm.com:5001/chaincode"

data1 = "{\"jsonrpc\": \"2.0\",\n\"method\": " #method type here
data2 = "\n\"params\": {\n\"type\": 1,\n\"chaincodeID\":{\n\"name\": \"12decb80905604f0f7ea238c51e8cd4c7496640865ba48b8798d8c35eb9be86b206c7119c925aa115513941a5c10dc9af34e6530ced4e673db2376553aa156ce\"\n},\n\"ctorMsg\":{\n\"function\": " #params.ctorMsg.function args
data3 = "\n\"args\": [" # params.ctorMsg.args here
data4 = "]\n},\n\"secureContext\": \"user_type1_1\"\n},\n\"id\": 0\n}"

#methodSel = 0

root = Tk()
methodSel = IntVar()

def radioSelect1():
    global methodSel
    if(methodSel.get() == 0):
        print("QUERYING")
    elif(methodSel.get() == 1):
        print("INVOKING")




def buttonClick():

    global methodSel
    # method = methodText.get("1.0", 'end-1c')
    function = functionText.get("1.0", 'end-1c')
    arguments = argumentsText.get("1.0", 'end-1c')

    print("MethodSel is: " + str(methodSel))

    if methodSel.get() == 0:
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
    print("Data Entered: " + str(data))
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

    print("Response: " + str(response))

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

    textArea2.config(state = NORMAL)
    textArea2.delete("1.0", END)
    textArea2.insert(END, response)
    textArea2.config(state = DISABLED)

    print("Data entered: \n\n" + data)
    print("Response received: \n\t" + r.text)

textArea2 = Text( height = 10, width = 40)
textArea2.grid(row = 0, column = 0, columnspan = 2)
textArea2.insert(END, "This is where we will receive text from Bluemix")
textArea2.config(state = DISABLED, font =("Times New Roman", 35))

# label1 = Label( height = 2, text = "Method: ")
# label1.grid(row = 1, column = 0, sticky = W)

label2 = Label(  height = 2, text = "Function: ")
label2.grid(row =2, column = 0, sticky = W)

label3 = Label(  height = 2, text = "Arguments: ")
label3.grid(row = 3, column = 0, sticky = W)

button1 = Button( command = buttonClick)
button1.grid(row = 4, column = 0, sticky = W)
button1.config(text = "Send Command", height = 1, width = 20)

# methodText = Text( height = 1, width = 20)
# methodText.grid(row = 1, column = 1)

functionText = Text( height = 1, width = 20)
functionText.grid(row = 2, column = 1, sticky = W)

argumentsText = Text( height = 1, width = 20)
argumentsText.grid(row = 3, column = 1, sticky = W)

"""=------------ BELOW Radio button Section -------------="""

queryButton = Radiobutton( font =("Ubuntu", 35), text = "Query", variable =methodSel, value = 0)
invokeButton = Radiobutton( font =("Ubuntu", 35), text = "Invoke", variable =methodSel, value = 1)

queryButton.config(borderwidth = 10, indicatoron=False, relief = RAISED, width = 10)
queryButton.grid(row = 1, column = 0, sticky = W)

invokeButton.config(borderwidth=10, indicatoron=False, relief = RAISED, width = 10)
invokeButton.grid(row = 1, column = 1, sticky = W)

"""=-----------------------------------------------------="""

QfunctionLabel = Label(text = "Query Function List:")
QfunctionLabel.grid(row = 5, column = 0, sticky = W)
QfunctionLabel.config(font = ("Arial", 14,"underline"))

QfunctionListLabel = Label(text = "read\nquery_functions\ninvoke_functions\nopen_trades\nview_my_assets")
QfunctionListLabel.grid(row=6, column=0, sticky = NW)
QfunctionListLabel.config(justify = LEFT)


IfunctionLabel = Label(text = "Invoke Function List:")
IfunctionLabel.grid(row=5, column = 1, sticky = W)
IfunctionLabel.config(font = ("Arial", 14,"underline"))

IfunctionListLabel = Label(text = "write\ndelete\ninit_energy\nset_owner\nopen_trade\nperform_trade\nremove_trade")
IfunctionListLabel.grid(row=6, column=1, sticky = NW)
IfunctionListLabel.config(justify = LEFT)


# textArea1 = Text( height = 20, width = 50)
# textArea1.grid(row = 0, column = 0)
# textArea1.insert(END, "This is Json sent to Bluemix")
# textArea1.config(text = "This is where text will be sent to Bluemix")


root.title("Bluemix Blockchain-Interface")
#root.geometry("1300x900")
root.mainloop()