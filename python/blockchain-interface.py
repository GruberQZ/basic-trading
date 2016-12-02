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
FuncSel = StringVar()

def parse(event):
    buttonClick()

def avoidIssue():
    qFunctionButtons = ['query_functions', 'invoke_functions', 'open_trades', 'view_my_assets', 'read']
    iFunctionButtons = ['init_energy', 'set_owner', 'open_trade', 'perform_trade', 'remove_trade',
                        'write', 'delete']
    if FuncSel.get() in qFunctionButtons and methodSel.get() == 1:
        iFuncIniTrade.select()
    if FuncSel.get() in iFunctionButtons and methodSel.get() == 0:
        qFuncOpenTrade.select()


def radioSelect1():
    global methodSel
    global iFuncWrite, iFuncSetOwn, iFuncRemTrade, iFuncPerfTrade, iFuncOpTrade, iFuncIniTrade, iFuncDel
    global qFuncQuerFunc, qFuncInvFunc, qFuncOpenTrade, qFuncViewMyAss, qFuncRead

    if(methodSel.get() == 0):
        print("QUERY")
        iFuncWrite.config(state = "disabled")
        iFuncSetOwn.config(state = "disabled")
        iFuncRemTrade.config(state = "disabled")
        iFuncPerfTrade.config(state = "disabled")
        iFuncOpTrade.config(state = "disabled")
        iFuncIniTrade.config(state = "disabled")
        iFuncDel.config(state = "disabled")

        qFuncQuerFunc.config(state = "normal")
        qFuncInvFunc.config(state = "normal")
        qFuncOpenTrade.config(state = "normal")
        qFuncViewMyAss.config(state = "normal")
        qFuncRead.config(state = "normal")

    elif(methodSel.get() == 1):
        print("INVOKE")
        qFuncQuerFunc.config(state = "disabled")
        qFuncInvFunc.config(state = "disabled")
        qFuncOpenTrade.config(state = "disabled")
        qFuncViewMyAss.config(state = "disabled")
        qFuncRead.config(state = "disabled")

        iFuncWrite.config(state = "normal")
        iFuncSetOwn.config(state = "normal")
        iFuncRemTrade.config(state = "normal")
        iFuncPerfTrade.config(state = "normal")
        iFuncOpTrade.config(state = "normal")
        iFuncIniTrade.config(state = "normal")
        iFuncDel.config(state = "normal")

    avoidIssue()


def buttonClick():

    global methodSel
    # method = methodText.get("1.0", 'end-1c')
    arguments = argumentsText.get("1.0", 'end-1c')

    print("MethodSel is: " + str(methodSel))

    if methodSel.get() == 0:
        method = "query"
    else:
        method = "invoke"
    # qFunctionButtons = ['query_functions', 'invoke_functions', 'open_trades', 'view_my_assets', 'read']

    function = FuncSel.get()

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
        print("Time to hang it up, let Kostas know we weren't ready :'(")
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
    response = response.replace('message:', '')
    response = response.replace('result:', '')

    textArea2.config(state = NORMAL)
    textArea2.delete("1.0", END)
    textArea2.insert(END, response)
    textArea2.config(state = DISABLED)

    print("Data entered: \n\n" + data)
    print("Response received: \n\t" + response)

textArea2 = Text( height = 10, width = 110)
textArea2.grid(row = 0, column = 0, columnspan = 2, padx = 10, ipadx = 10, pady = 20)
textArea2.insert(END, "This is where we will receive text from Bluemix")
textArea2.config(state = DISABLED, font =("Times New Roman", 22))

# label1 = Label( height = 2, text = "Method: ")
# label1.grid(row = 1, column = 0, sticky = W)

label2 = Label(  height = 2, text = "Function: ")
label2.grid(row =2, column = 0, padx = 10)
label2.config(font = ("Arial", 20))

label3 = Label(  height = 2, text = "Arguments: ")
label3.grid(row = 9, column = 0, padx = 10)
label3.config(font = ("Arial", 20))

button1 = Button( command = buttonClick)
button1.grid(row = 3, column = 0, padx = 10, columnspan = 2)
button1.config(text = "Send Command", height = 1, width = 20, font = ("Arial", 20, "bold"), bg ="#CC0000", fg = "white")

# methodText = Text( height = 1, width = 20)
# methodText.grid(row = 1, column = 1)

functionLabel = Label(text = "Functions List")
functionLabel.config(font=("Arial", 18, "underline"))
functionLabel.grid(row = 2, columnspan = 2)

iFuncWrite = Radiobutton(text = 'write', variable = FuncSel, value = 'write')
iFuncWrite.grid(row=7, column=1)
iFuncWrite.config(font=("Arial", 16))

iFuncDel = Radiobutton(text = 'delete', variable = FuncSel, value = 'delete')
iFuncDel.grid(row=8, column=1)
iFuncDel.config(font=("Arial", 16))

iFuncSetOwn = Radiobutton(text = 'set_owner', variable = FuncSel, value = 'set_owner')
iFuncSetOwn.grid(row=3, column=1)
iFuncSetOwn.config(font=("Arial", 16))

iFuncOpTrade = Radiobutton(text = 'open_trade', variable = FuncSel, value = 'open_trade')
iFuncOpTrade.grid(row=4, column=1)
iFuncOpTrade.config(font=("Arial", 16))

iFuncPerfTrade = Radiobutton(text = 'perform_trade', variable = FuncSel, value = 'perform_trade')
iFuncPerfTrade.grid(row=5, column=1)
iFuncPerfTrade.config(font=("Arial", 16))

iFuncRemTrade = Radiobutton(text = 'remove_trade', variable = FuncSel, value = 'remove_trade')
iFuncRemTrade.grid(row=6, column=1)
iFuncRemTrade.config(font=("Arial", 16))

iFuncIniTrade = Radiobutton(text = 'init_energy', variable = FuncSel, value = 'init_energy')
iFuncIniTrade.grid(row=2, column=1)
iFuncIniTrade.config(font=("Arial", 16))

qFunctionButtons = ['query_functions', 'invoke_functions', 'open_trades', 'view_my_assets', 'read']
qFuncQuerFunc = Radiobutton(text=qFunctionButtons[0], variable=FuncSel, value=qFunctionButtons[0])
qFuncQuerFunc.grid(row=4, column=0)
qFuncQuerFunc.config(font=("Arial", 16))

qFuncInvFunc = Radiobutton(text=qFunctionButtons[1], variable=FuncSel, value=qFunctionButtons[1])
qFuncInvFunc.grid(row=5, column=0)
qFuncInvFunc.config(font=("Arial", 16))

qFuncOpenTrade = Radiobutton(text=qFunctionButtons[2], variable=FuncSel, value=qFunctionButtons[2])
qFuncOpenTrade.grid(row=2, column=0)
qFuncOpenTrade.config(font=("Arial", 16))

qFuncViewMyAss = Radiobutton(text=qFunctionButtons[3], variable=FuncSel, value=qFunctionButtons[3])
qFuncViewMyAss.grid(row=3, column=0)
qFuncViewMyAss.config(font=("Arial", 16))

qFuncRead = Radiobutton(text=qFunctionButtons[4], variable=FuncSel, value=qFunctionButtons[4])
qFuncRead.grid(row=6, column=0)
qFuncRead.config(font=("Arial", 16))


# rowStart = 2
# for i in qFunctionButtons:
#     b = Radiobutton(text=i, variable=iFuncSel, value=i)
#     b.grid(row=rowStart, column=0)
#     b.config(font=("Arial", 20))
#     rowStart += 1

argumentsText = Text( height = 1, width = 20)
argumentsText.grid(row = 9, column = 1, ipadx = 10)
argumentsText.config(font = ("Arial", 20))
argumentsText.bind('<Return>', parse)


"""=------------ BELOW Radio button Section -------------="""

queryButton = Radiobutton( font =("Ubuntu", 22), text = "Query", command = radioSelect1,variable =methodSel, value = 0)
invokeButton = Radiobutton( font =("Ubuntu", 22), text = "Invoke", command = radioSelect1, variable =methodSel, value = 1)

queryButton.config(borderwidth = 10, indicatoron=False, relief = RAISED, width = 10)
queryButton.grid(row = 1, column = 0, padx = 10)

invokeButton.config(borderwidth=10, indicatoron=False, relief = RAISED, width = 10)
invokeButton.grid(row = 1, column = 1, ipadx = 10)

"""=-----------------------------------------------------="""

QfunctionLabel = Label(text = "Query Function List:")
QfunctionLabel.grid(row = 11, column = 0, padx = 10)
QfunctionLabel.config(font = ("Arial", 20,"underline"))

QfunctionListLabel = Label(text = "open_trades\nview_my_assets: \t\"Owner\"\nquery_functions\ninvoke_functions\nread: \t\t\"Variable\"")
QfunctionListLabel.grid(row=12, column=0, padx = 10, stick = N)
QfunctionListLabel.config(justify = LEFT, font=("Arial", 14))


IfunctionLabel = Label(text = "Invoke Function List:")
IfunctionLabel.grid(row=11, column = 1, ipadx = 10)
IfunctionLabel.config(font = ("Arial", 20,"underline"))

IfunctionListLabel = Label(text = "init_energy: \t\"UniqueID,Amount,\n\t\tEnergyPrice,Owner\"\nset_owner: \t\"Asset,Owner\"\n"
                                  "open_trade: \t\"Owner,Asset\"\nperform_trade: \t\"Asset,Owner\"\n"
                                  "remove_trade: \t\"Owner,Asset\"\nwrite: \t\t\"Variable,Value\"\ndelete: \t\t\"Variable\"\n")
IfunctionListLabel.grid(row=12, column=1, ipadx = 10)
IfunctionListLabel.config(justify = LEFT, font=("Arial", 14))

# argLabel = Label(text = "Arguments Available = \"Example\"")
# argLabel.grid(row=11, ipadx = 11, columnspan = 2, pady = 5)
# argLabel.config(justify = CENTER, font=("Arial", 18, "underline"))
#
# argList1Label = Label(text = "Asset Name = \"Energy1\"\nClient Name = \"Bob\"\nGasPrice = \"25\"")
# argList1Label.grid(row=12, column=0, padx = 10)
# argList1Label.config(justify = LEFT, font=("Arial", 14))
#
# argList2Label = Label(text = "EnergyPrice = \"100\"\nVariable = \"ece\"\nValue to write to Variable = \"200\"")
# argList2Label.grid(row=12, column=1, ipadx = 10)
# argList2Label.config(justify = LEFT, font=("Arial", 14))

radioSelect1()
qFuncOpenTrade.select()
root.rowconfigure(13, weight = 1)
root.rowconfigure(14, weight = 1)
root.columnconfigure(0, weight = 1)
root.columnconfigure(1, weight = 1)
root.title("Bluemix Blockchain-Interface")
root.mainloop()