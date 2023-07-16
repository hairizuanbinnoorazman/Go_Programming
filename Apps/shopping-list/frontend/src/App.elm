module App exposing (..)

import Browser
import Browser.Navigation as Nav
import Html exposing (..)
import Html.Attributes exposing (..)
import Http as Http
import Json.Decode as Decode exposing (Decoder, bool, decodeString, float, int, list, null, string)
import Json.Decode.Pipeline as Pipeline
import Json.Encode as Encode
import Material.Button as Button
import Material.Dialog as Dialog
import Material.LayoutGrid as LayoutGrid
import Material.List as List
import Material.List.Item as ListItem
import Material.TextField as TextField
import Url
import Url.Parser as Url exposing ((</>), (<?>), Parser)


type alias Flags =
    { serverEndpoint : String
    , ingressPath : String
    }


main : Program Flags Model Msg
main =
    Browser.application
        { init = init
        , view = view
        , update = update
        , subscriptions = subscriptions
        , onUrlChange = UrlChanged
        , onUrlRequest = LinkClicked
        }


type alias Model =
    { key : Nav.Key
    , url : Url.Url
    , page : Page
    , serverSettings : Flags
    , dialogOpen : Bool
    , dialogItem : String
    , shoppingItems : List ShoppingItem
    }


type Page
    = Index


urlToPage : String -> Url.Url -> Page
urlToPage ingressPath url =
    url
        |> Url.parse (urlParser ingressPath)
        |> Maybe.withDefault Index


urlParser : String -> Parser (Page -> a) a
urlParser ingressPath =
    -- We try to match one of the following URLs
    Url.oneOf
        [ Url.map Index (Url.s ingressPath </> Url.top)
        ]


init : Flags -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
    let
        initialAppState =
            { key = key
            , url = url
            , page = urlToPage flags.ingressPath url
            , serverSettings = flags
            , dialogOpen = False
            , dialogItem = ""
            , shoppingItems = []
            }
    in
    ( initialAppState, Cmd.batch [ Nav.pushUrl key (flags.ingressPath ++ "/"), apiListShoppingItems flags.serverEndpoint ] )


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | Clicked
    | Closed
    | ValueChanged String
    | GetShoppingItemListResponse (Result Http.Error ShoppingItemListResp)
    | SingleItemClicked ShoppingItem
    | ToggleShoppingItem (Result Http.Error ())
    | DeleteShoppingItemResponse (Result Http.Error ())
    | CreateShoppingItemResponse (Result Http.Error ())
    | ClearClicked


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UrlChanged url ->
            let
                indexURL =
                    { url | path = model.serverSettings.ingressPath ++ "/" }
            in
            ( model, Cmd.none )

        LinkClicked urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        Clicked ->
            ( { model | dialogOpen = True, dialogItem = "" }, Cmd.none )

        Closed ->
            ( { model | dialogOpen = False, dialogItem = "" }, Cmd.batch [ apiCreateShoppingItem model.serverSettings.serverEndpoint model.dialogItem ] )

        ValueChanged zzz ->
            ( { model | dialogItem = zzz }, Cmd.none )

        GetShoppingItemListResponse result ->
            case result of
                Ok zzz ->
                    ( { model | shoppingItems = zzz.items }, Cmd.none )

                Err zzz ->
                    ( model, Cmd.none )

        SingleItemClicked item ->
            let
                updatedItems =
                    List.map (toggleStatus item.id) model.shoppingItems
            in
            ( { model | shoppingItems = updatedItems }, Cmd.batch [ apiUpdateShoppingItem model.serverSettings.serverEndpoint item.id (toggleStatusString item.status) ] )

        ToggleShoppingItem result ->
            ( model, Cmd.none )

        DeleteShoppingItemResponse result ->
            ( model, Cmd.none )

        CreateShoppingItemResponse result ->
            ( model, Cmd.batch [ apiListShoppingItems model.serverSettings.serverEndpoint ] )

        ClearClicked ->
            let
                incartItems =
                    List.filter (filterItem "in-cart") model.shoppingItems

                activeItems =
                    List.filter (filterItem "active") model.shoppingItems

                deleteCommands =
                    deleteCmds model.serverSettings.serverEndpoint incartItems
            in
            ( { model | shoppingItems = activeItems }, Cmd.batch deleteCommands )


deleteCmds : String -> List ShoppingItem -> List (Cmd Msg)
deleteCmds mgrUrl items =
    List.map (deleteCmd mgrUrl) items


deleteCmd : String -> ShoppingItem -> Cmd Msg
deleteCmd mgrURL item =
    apiDeleteShoppingItem mgrURL item.id


filterItem : String -> ShoppingItem -> Bool
filterItem status item =
    if item.status == status then
        True

    else
        False


toggleStatusString : String -> String
toggleStatusString status =
    if status == "active" then
        "in-cart"

    else
        "active"


toggleStatus : String -> ShoppingItem -> ShoppingItem
toggleStatus id shoppingItem =
    if shoppingItem.id == id then
        if shoppingItem.status == "active" then
            { shoppingItem | status = "in-cart" }

        else
            { shoppingItem | status = "active" }

    else
        shoppingItem


errorToString : Http.Error -> String
errorToString error =
    case error of
        Http.BadUrl url ->
            "The URL " ++ url ++ " was invalid"

        Http.Timeout ->
            "Unable to reach the server, try again"

        Http.NetworkError ->
            "Unable to reach the server, check your network connection"

        Http.BadStatus 500 ->
            "The server had a problem, try again later"

        Http.BadStatus 400 ->
            "Verify your information and try again"

        Http.BadStatus _ ->
            "Unknown error"

        Http.BadBody errorMessage ->
            errorMessage


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


view : Model -> Browser.Document Msg
view model =
    { title = "Shopping List"
    , body =
        [ Dialog.alert
            (Dialog.config
                |> Dialog.setOpen model.dialogOpen
                |> Dialog.setOnClose Closed
            )
            { content =
                [ TextField.filled
                    (TextField.config
                        |> TextField.setLabel (Just "Item at Shopping List")
                        |> TextField.setValue (Just model.dialogItem)
                        |> TextField.setOnInput ValueChanged
                    )
                ]
            , actions =
                [ Button.text
                    (Button.config
                        |> Button.setOnClick Closed
                        |> Button.setAttributes [ Dialog.defaultAction ]
                    )
                    "Submit"
                ]
            }
        , LayoutGrid.layoutGrid [ LayoutGrid.span12 ]
            [ h1 [] [ text "Shopping List" ]
            , Button.text
                (Button.config |> Button.setOnClick Clicked)
                "Add item"
            , Button.text
                (Button.config |> Button.setOnClick ClearClicked)
                "Clear in-cart items"
            , LayoutGrid.inner []
                [ LayoutGrid.cell [ LayoutGrid.span12 ]
                    [ List.list List.config
                        (ListItem.listItem ListItem.config
                            [ h2 [] [ text "Shopping Items" ] ]
                        )
                        (List.map singleShoppingItemRow model.shoppingItems)
                    ]
                ]
            ]
        ]
    }


singleShoppingItemRow : ShoppingItem -> ListItem.ListItem Msg
singleShoppingItemRow item =
    if item.status == "active" then
        ListItem.listItem (ListItem.config |> ListItem.setOnClick (SingleItemClicked item)) [ p [ id item.id ] [ text item.name ] ]

    else
        ListItem.listItem (ListItem.config |> ListItem.setOnClick (SingleItemClicked item) |> ListItem.setAttributes [ style "background-color" "yellow" ]) [ p [ id item.id ] [ text item.name ] ]


type alias ShoppingItem =
    { id : String
    , name : String
    , created : String
    , status : String
    }


shoppingItemDecoder : Decode.Decoder ShoppingItem
shoppingItemDecoder =
    Decode.succeed ShoppingItem
        |> Pipeline.required "ID" string
        |> Pipeline.required "name" string
        |> Pipeline.required "created" string
        |> Pipeline.required "status" string


type alias ShoppingItemListResp =
    { items : List ShoppingItem
    }


shoppingItemListDecoder : Decode.Decoder ShoppingItemListResp
shoppingItemListDecoder =
    Decode.succeed ShoppingItemListResp
        |> Pipeline.required "items" (Decode.list shoppingItemDecoder)


apiListShoppingItems : String -> Cmd Msg
apiListShoppingItems mgrURL =
    let
        url =
            mgrURL ++ "/api/shopping-list/v1/item"
    in
    Http.request
        { body = Http.emptyBody
        , method = "GET"
        , url = url
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        , expect = Http.expectJson GetShoppingItemListResponse shoppingItemListDecoder
        }


apiCreateShoppingItem : String -> String -> Cmd Msg
apiCreateShoppingItem mgrURL name =
    let
        url =
            mgrURL ++ "/api/shopping-list/v1/item"

        statusBody =
            Http.jsonBody <|
                Encode.object
                    [ ( "name", Encode.string name ) ]
    in
    Http.request
        { body = statusBody
        , method = "POST"
        , url = url
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        , expect = Http.expectWhatever CreateShoppingItemResponse
        }


apiUpdateShoppingItem : String -> String -> String -> Cmd Msg
apiUpdateShoppingItem mgrURL shoppingItemID status =
    let
        url =
            mgrURL ++ "/api/shopping-list/v1/item/" ++ shoppingItemID

        statusBody =
            Http.jsonBody <|
                Encode.object
                    [ ( "status", Encode.string status ) ]
    in
    Http.request
        { body = statusBody
        , method = "PATCH"
        , url = url
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        , expect = Http.expectWhatever ToggleShoppingItem
        }


apiDeleteShoppingItem : String -> String -> Cmd Msg
apiDeleteShoppingItem mgrURL shoppingItemID =
    let
        url =
            mgrURL ++ "/api/shopping-list/v1/item/" ++ shoppingItemID
    in
    Http.request
        { body = Http.emptyBody
        , method = "DELETE"
        , url = url
        , headers = []
        , timeout = Nothing
        , tracker = Nothing
        , expect = Http.expectWhatever DeleteShoppingItemResponse
        }
