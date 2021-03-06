package types

import (
    "net/http"
)

type Auth struct {
    
    Email     string
    Password  string
    Token     string
    UserId   int
    //Auth remove this line for disable generator functionality
}

func (auth *Auth) Validate()  {
    //Validate remove this line for disable generator functionality
}

type AuthFilter struct {
    model Auth
    //AuthFilter remove this line for disable generator functionality

    AbstractFilter
}

func GetAuthFilter(request *http.Request, functionType string) AuthFilter {

    var filter AuthFilter

    filter.request = request
    //filter.TestFilter, _ = strconv.Atoi(request.FormValue("TestFilter"))

    //GetAuthFilter remove this line for disable generator functionality

    ReadJSON(request, &filter.model)

    filter.AbstractFilter = GetAbstractFilter(request, functionType)

    return  filter
}


func (filter *AuthFilter) GetAuthModel() Auth {

    filter.model.Validate()

    return  filter.model
}

func (filter *AuthFilter) SetAuthModel(typeModel Auth) {

    filter.model = typeModel
}
