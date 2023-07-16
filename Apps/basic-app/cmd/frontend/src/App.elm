module App exposing (..)

import Bootstrap.CDN as CDN
import Bootstrap.Grid as Grid
import Bootstrap.Navbar as Navbar
import Browser
import Browser.Navigation as Nav
import Http
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
    , navbarState : Navbar.State
    }


type Page
    = Index
    | Login


urlParser : Parser (Page -> a) a
urlParser =
    -- We try to match one of the following URLs
    Url.oneOf
        [ Url.map Index Url.top
        , Url.map Login (Url.s "login")
        ]


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | NavbarMsg Navbar.State


init : Flags -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init flags url key =
    let
        ( navbarState, navbarCmd ) =
            Navbar.initialState NavbarMsg

        initialAppState =
            { key = key
            , url = url
            , navbarState = navbarState
            }
    in
    ( initialAppState, Cmd.none )


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
    Sub.batch
        []


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        LinkClicked urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Nav.load href )

        UrlChanged url ->
            ( model, Cmd.none )

        NavbarMsg state ->
            ( { model | navbarState = state }, Cmd.none )


view : Model -> Browser.Document Msg
view model =
    { title = "Basic Application"
    , body =
        [ Grid.container []
            [ CDN.stylesheet
            , Grid.row []
                [ Grid.col [] []
                ]
            ]
        ]
    }
