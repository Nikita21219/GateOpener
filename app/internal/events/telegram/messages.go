package telegram

const msgHelp = `
/open_entry - открыть шлагбаум на въезд

/open_exit - открыть шлагбаум на выезд

/opening_mode - открыть оба шлагбаума на 5 минут

/open_exit_mode - открыть шлагбаум на выезд на 30 секунд

/opening_mode_stop - прекратить открывать шлагбаумы
`

const msgStart = `
Я могу помочь открыть шлагбаум.
Для более подробной информации введи /help
`

const MsgUnknownCommand = `Я тебя не понимаю`
const msgCantGateOpen = `Произошла ошибка при работе со шлагбаумом`
const msgGateOpened = `Шлагбаум открыт`
const msgGateOpeningModeActivated = `Шлагбаумы открыты на 5 минут`
const msgGateExitOpeningModeActivated = `Шлагбаум на выезд открыт на 30 секунд`
const msgGateOpeningModeDeactivated = `Режим открытия шлагбаумов деактивирован`
const MsgNotAllowedControl = `Вы не являетесь администратором`
const MsgOpeningModeStopped = `Режим открытия остановлен`
