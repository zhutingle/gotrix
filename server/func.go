package server
func init() {
	StaticResource["gotrix.xml"] = "PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0iVVRGLTgiID8+CjxmdW5jcz4KCiAgICA8ZnVuYyBpZD0iMCIgbmFtZT0iR2V0UGFzcyIgZGVzPSLku45UT0tFTuS4reWPlnBhc3MiIHByaXZhdGU9InRydWUiPgogICAgICAgIDxzdHJpbmcgbmFtZT0idG9rZW4iIG11c3Q9InRydWUiIGxlbj0iNjQtNjQiLz4KCiAgICAgICAgPGpvYiByZXN1bHQ9InVzZXJJbmZvIiB0eXBlPSJzaW5nbGUiPnNlbGVjdCAqIGZyb20gdXNlciB3aGVyZSB0b2tlbiA9ICR7dG9rZW59PC9qb2I+CiAgICAgICAgPGpvYj5FcSgke3VzZXJJbmZvfSxudWxsLCfor6XnlKjmiLfkuI3lrZjlnKgnKTwvam9iPgogICAgICAgIDxqb2I+SmdldCgke3VzZXJJbmZvfSwncGFzcycsJ2tleScpPC9qb2I+CiAgICA8L2Z1bmM+CgogICAgPGZ1bmMgaWQ9IjEiIG5hbWU9IkdldFNlc3Npb24iIGRlcz0i5qC55o2uc2Vzc2lvbueahOWAvO+8jOiOt+WPluino+WvhuivpXNlc3Npb27nmoTlr4bnoIEiIHByaXZhdGU9InRydWUiPgogICAgICAgIDxzdHJpbmcgbmFtZT0idG9rZW4iIG11c3Q9InRydWUiIGxlbj0iMjAtNjAiLz4KCiAgICAgICAgPGpvYiByZXN1bHQ9InNlc3Npb24iPkdldFNlc3Npb24oJHt0b2tlbn0pPC9qb2I+CiAgICAgICAgPGpvYj5FcSgke3Nlc3Npb259LG51bGwsJ+eUqOaIt+S8muivneW3sui/h+acn++8jOivt+mHjeaWsOeZu+mZhicpPC9qb2I+CiAgICAgICAgPGpvYj5KZ2V0KCR7c2Vzc2lvbn0sJ2lkJywncGFzcycpPC9qb2I+CiAgICA8L2Z1bmM+CgogICAgPGZ1bmMgaWQ9IjEwIiBuYW1lPSJHZXRTYWx0IiBkZXM9IuiOt+WPluWvhueggeeahOebkOWAvCI+CgogICAgICAgIDxqb2IgdHlwZT0ic2luZ2xlIj5zZWxlY3Qgc2FsdCBmcm9tIHVzZXIgd2hlcmUgdG9rZW4gPSAke3Rva2VufTwvam9iPgogICAgICAgIDxqb2I+SmdldCgkezF9LCdzYWx0Jyk8L2pvYj4KICAgIDwvZnVuYz4KCiAgICA8ZnVuYyBpZD0iMTEiIG5hbWU9IkxvZ2luSW4iIGRlcz0i55m76ZmGIj4KICAgICAgICA8c3RyaW5nIG5hbWU9InRva2VuIiBtdXN0PSJ0cnVlIiBsZW49IjY0LTY0Ii8+CiAgICAgICAgPHN0cmluZyBuYW1lPSJ4IiBtdXN0PSJ0cnVlIiBsZW49IjIwLTUwIi8+CiAgICAgICAgPHN0cmluZyBuYW1lPSJ5IiBtdXN0PSJ0cnVlIiBsZW49IjIwLTUwIi8+CgogICAgICAgIDxqb2IgcmVzdWx0PSJ1c2VySW5mbyIgdHlwZT0ic2luZ2xlIj5zZWxlY3QgaWQsc2Vzc2lvbiBvbGRTZXNzaW9uIGZyb20gdXNlciB3aGVyZSB0b2tlbiA9ICR7dG9rZW59PC9qb2I+CiAgICAgICAgPGpvYiByZXN1bHQ9Im9sZFNlc3Npb24iPkpnZXQoJHt1c2VySW5mb30sJ29sZFNlc3Npb24nKTwvam9iPgogICAgICAgIDxqb2IgdGVzdD0iTmVxKCR7b2xkU2Vzc2lvbn0sbnVsbCkiPkRlbFNlc3Npb24oJHtvbGRTZXNzaW9ufSk8L2pvYj4KCiAgICAgICAgPGpvYiByZXN1bHQ9ImlkIj5KZ2V0KCR7dXNlckluZm99LCdpZCcpPC9qb2I+CiAgICAgICAgPGpvYiByZXN1bHQ9ImxvZ2luIj5Mb2dpbkluKCR7eH0sJHt5fSwke2lkfSk8L2pvYj4KICAgICAgICA8am9iIHJlc3VsdD0ic2Vzc2lvbiI+SmdldCgke2xvZ2lufSwnc2Vzc2lvbicpPC9qb2I+CiAgICAgICAgPGpvYiByZXN1bHQ9InBhc3MiPkpnZXQoJHtsb2dpbn0sJ3Bhc3MnKTwvam9iPgogICAgICAgIDxqb2IgcmVzdWx0PSJqc29uIj5Kc2V0KCdpZCcsJHtpZH0sJ3Nlc3Npb24nLCR7c2Vzc2lvbn0sJ3Bhc3MnLCR7cGFzc30pPC9qb2I+CiAgICAgICAgPGpvYj51cGRhdGUgdXNlciBzZXQgc2Vzc2lvbiA9ICR7c2Vzc2lvbn0gd2hlcmUgaWQgPSAke2lkfTwvam9iPgogICAgICAgIDxqb2I+U2V0U2Vzc2lvbigke3Nlc3Npb259LCR7anNvbn0pPC9qb2I+CiAgICAgICAgPGpvYj5KZ2V0KCR7bG9naW59LCdzZXNzaW9uJywneCcsJ3knKTwvam9iPgogICAgPC9mdW5jPgoKICAgIDxmdW5jIGlkPSIxMiIgbmFtZT0iTG9naW5PdXQiIGRlcz0i55m75Ye6Ij4KICAgICAgICA8am9iPkRlbFNlc3Npb24oJHt0b2tlbn0pPC9qb2I+CiAgICA8L2Z1bmM+CgogICAgPGZ1bmMgaWQ9IjEzIiBuYW1lPSJyZWdpc3RlciIgZGVzPSLnlKhUT0tFTuadpeazqOWGjCIgc2VsZj0idHJ1ZSI+CiAgICAgICAgPHN0cmluZyBuYW1lPSJ0b2tlbiIgbXVzdD0idHJ1ZSIgbGVuPSI2NC02NCIvPgogICAgICAgIDxzdHJpbmcgbmFtZT0icGFzcyIgbXVzdD0idHJ1ZSIgbGVuPSIyMC0zMiIvPgogICAgICAgIDxzdHJpbmcgbmFtZT0ic2FsdCIgbXVzdD0idHJ1ZSIgbGVuPSIyMC0zMiIvPgogICAgICAgIDxzdHJpbmcgbmFtZT0ia2V5IiBtdXN0PSJ0cnVlIiBsZW49IjIwLTMyIi8+CiAgICAgICAgPHN0cmluZyBuYW1lPSJ0ZWwiIGxlbj0iMC0yMCIvPgogICAgICAgIDxzdHJpbmcgbmFtZT0ibmljayIgbGVuPSIwLTY0Ii8+CiAgICAgICAgPHN0cmluZyBuYW1lPSJwb3J0cmFpdCIgbGVuPSIwLTEyOCIvPgogICAgICAgIDxpbnQgbmFtZT0icGFyZW50Ii8+CgogICAgICAgIDxqb2IgdHlwZT0ic2luZ2xlIj5zZWxlY3QgY291bnQoKikgY291bnQgZnJvbSB1c2VyIHdoZXJlIHRva2VuID0gJHt0b2tlbn08L2pvYj4KICAgICAgICA8am9iPkpnZXQoJHsxfSwnY291bnQnKTwvam9iPgogICAgICAgIDxqb2I+RXEoJHsyfSwxLCfotKblj7flt7LlrZjlnKjvvIzlj6/ku6Xnm7TmjqXnmbvpmYYnKTwvam9iPgogICAgICAgIDxqb2I+aW5zZXJ0IGludG8gdXNlciB2YWx1ZXMobnVsbCwke3RlbH0sJHtuaWNrfSwke3BvcnRyYWl0fSwke3Rva2VufSwke3Bhc3N9LCR7c2FsdH0sJHtrZXl9LCR7cGFyZW50fSxudWxsKTwvam9iPgoKICAgIDwvZnVuYz4KCjwvZnVuY3M+Cg=="
}
