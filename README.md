# Taavi

## Tööriistad

- [Go (>= v1.18)](https://go.dev/)
- [Buf](https://buf.build/)
- [Docker (valikuline)](https://www.docker.com/)

## Seadistus

Käivita kasutas `proto` käsk `buf generate`. See genereerib vajalikud protobuf failid. Selle käsu jooksutamine on vajalik iga kord, kui protobuf faile on muudetud.

## Struktuur

### Zulipi bot

Kaustas (*package*) `zlp` on Zulipi integratsioon.

Boti loomiseks sobib järgmine koodilõik:
```go
func main() {
  rc, _ := zlp.LoadRC(".zuliprc")
  bot := zlp.NewBot(rc)
  bot.Init()
}
```

ZulipRC faili saab alla laadida Zulipi boti menüüst. Seda Git-i ei tohi panna!

### CRON tööd

Praegu käivitatakse tööd julmalt Taavi käivitumisel. Tulevikus peaks lisatööd käivitama nende lisamisel ja mingeid uuendusi tegema keskööl.

