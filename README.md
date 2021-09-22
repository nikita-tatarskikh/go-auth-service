# go-auth-service
https://www.notion.so/Test-task-Junior-BackDev-215fcddafff2425a8ca7e515e21527e7
- started at 127.0.0.1
- 127.0.0.1:8080/sign-up - принимает post запрос, content-type: x-www-form-urlencoded, ожидаемый key:guid
возвращает пару токенов access и refresh.
пример body:
{
    "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyZWZlcnNoX3V1aWQiOiIxYjkyYjAxYS1jZGI4LTRiY2QtYmJmNy00NWM2OWE1MzIzYjUiLCJ1c2VyX2lkIjoiYmNmMThkMjktNGM2Zi00MTJmLTk0MDAtN2U3NmQzMzI4ODc3In0.OnK8sxqnTXqwnvFU4zaGQCTjYH33KWRZ9lk4IgQWfcKDH9KuPoSxMMNaKg9Q28mfLDQdTDKqPCOGdNOMMxuBKQ",
    "refreshToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyZWZlcnNoX3V1aWQiOiJhMDk5YThmNS0zZmNlLTRmM2EtOGEyYi0wZGNiYmRiOWUwNWQiLCJ1c2VyX2lkIjoiYmNmMThkMjktNGM2Zi00MTJmLTk0MDAtN2U3NmQzMzI4ODc3In0.guiSdCzE8YkS6HM9cardXv1CvoU5nV0oR7gowdOIe6QogTfuUD1g2s2lqxv77hy1k-n4FEMgHIogNi0rNymSow"
}
- 127.0.0.1:8080/refresh - принмает post запрос, content-type: application-json,
body: {
   "refreshToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyZWZlcnNoX3V1aWQiOiI0MmZiMmRhYi0xN2ZlLTQ2YjktYWI2OS1hMDM2ZTZjMGQ2NmQiLCJ1c2VyX2lkIjoiYmNmMThkMjktNGM2Zi00MTJmLTk0MDAtN2U3NmQzMzI4ODc3In0.8cnlnvtRXY2q4RkZlmPZkTps2BYweIrcDUbgWU9N6wn4rZUBd1MlgqRZdQbFm6K1WDXrXRo0bvKA0bX9ntuRpg"
}
возвращает:

пример body:
{
    "accessToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyZWZlcnNoX3V1aWQiOiJkNjA2ZmI5Zi0xMzY2LTQ4YWMtOTA4MS00M2MyY2VhNGY4YmIiLCJ1c2VyX2lkIjoiYmNmMThkMjktNGM2Zi00MTJmLTk0MDAtN2U3NmQzMzI4ODc3In0.cbXWPrambc_Rs_m7BCpiOhOlc06ECz4IavW4SQfc8ktorOCJy0bbn-TBZLqMVIT0VU7P57ImvjveXohG1xXDuQ",
    "refreshToken": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJyZWZlcnNoX3V1aWQiOiJiNmM5Njk4Yi05YTYyLTQ5NGItODZhNC1jMjA5ZTg4ZDQ4YjciLCJ1c2VyX2lkIjoiYmNmMThkMjktNGM2Zi00MTJmLTk0MDAtN2U3NmQzMzI4ODc3In0.trA4YM1b6jsoqfJMuKsPFNIjTkZNmKIoygs8B_vRPS1ydOAKtW2GfNt9lrQgluEkFkk-WTIsoczLmJXsg1OwiA"
}

