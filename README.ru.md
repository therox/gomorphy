<p align="center">
  <img src="assets/mikhalich.png" alt="Михалыч — талисман gomorphy" width="180">
</p>

<h1 align="center">gomorphy</h1>
[![Go Reference](https://pkg.go.dev/badge/github.com/therox/gomorphy.svg)](https://pkg.go.dev/github.com/therox/gomorphy)
Словарная библиотека для русской морфологии: склонение существительных,
прилагательных и ФИО, согласование с числительными, морфологический разбор
словоформы.

Реализовано: `Decline`, `DeclineAdj`, `Agree`, `GenderOf`, `PluralOf`,
`Parse`, `DeclineFullName`, `ToNominative`, `ParseFullName` — включая
прямые и обратные эвристики для несловарных отчеств и фамилий.

> English README: [README.md](README.md).

## Установка

```sh
go get github.com/therox/gomorphy
```

Зависимостей нет — только стандартная библиотека.

## Скачивание и сборка словаря

Полный словарь OpenCorpora распространяется под лицензией CC BY-SA и в репозиторий
не входит. Размер дампа — ~16 МБ bz2 / ~500 МБ распакованного XML. Дамп
обновляется на стороне OpenCorpora — актуальный список форматов на странице
[opencorpora.org/?page=downloads](https://opencorpora.org/?page=downloads).

```sh
# 1. Скачиваем дамп (bz2 ~16 МБ; есть и .zip ~27 МБ).
curl -L -o dict.opcorpora.xml.bz2 \
    https://opencorpora.org/files/export/dict/dict.opcorpora.xml.bz2

# 2. Распаковываем (~500 МБ XML).
bunzip2 dict.opcorpora.xml.bz2

# 3. Собираем компактный .bin для библиотеки.
go run ./cmd/builddict -in dict.opcorpora.xml -out dict.bin
```

Для быстрой проверки без скачивания дампа есть мини-словарь
`testdata/sample.xml` — на нём проходят все примеры из этого README:

```sh
go run ./cmd/builddict -in testdata/sample.xml -out dict.bin
```

## Инициализация словаря

Перед первым обращением к API словарь должен быть загружен в память — либо
явно через `Init`, либо автоматически из переменной окружения `GOMORPHY_DICT`:

```go
// Явная инициализация.
if err := gomorphy.Init("dict.bin"); err != nil {
    log.Fatal(err)
}

// Альтернатива: задать GOMORPHY_DICT=/path/to/dict.bin перед запуском —
// и первый же вызов API подгрузит словарь.
```

`Init` можно вызвать только один раз; повторный вызов вернёт ошибку.

## Использование

Все примеры ниже соответствуют ассертам из `*_test.go` в корне репозитория
и работают на мини-словаре `testdata/sample.xml`.

```go
package main

import (
    "fmt"
    "log"

    "github.com/therox/gomorphy"
)

func main() {
    if err := gomorphy.Init("dict.bin"); err != nil {
        log.Fatal(err)
    }

    // Существительное: «аппетит» в род.п. ед.ч. → «аппетита».
    word, _ := gomorphy.Decline("аппетит", gomorphy.Genitive, gomorphy.Singular)
    fmt.Println(word) // аппетита

    // Прилагательное: «красный» в вин.п. ед.ч., муж.р., одуш. → «красного».
    adj, _ := gomorphy.DeclineAdj("красный", gomorphy.Accusative,
        gomorphy.Singular, gomorphy.Masculine, true)
    fmt.Println(adj) // красного

    // Согласование с числительным: 1 яблоко / 2 яблока / 5 яблок / 12 яблок.
    for _, n := range []int{1, 2, 5, 12} {
        s, _ := gomorphy.Agree("яблоко", n)
        fmt.Println(n, s)
    }
    // 1 яблоко
    // 2 яблока
    // 5 яблок
    // 12 яблок

    // Род слова.
    g, _ := gomorphy.GenderOf("книга")
    fmt.Println(g == gomorphy.Feminine) // true

    // Множественное число (включая супплетивные пары).
    pl, _ := gomorphy.PluralOf("человек")
    fmt.Println(pl) // люди

    // Морфологический разбор: «стали» — род./дат./предл. ед.ч. + им.п. мн.ч.
    analyses, _ := gomorphy.Parse("стали")
    for _, a := range analyses {
        fmt.Printf("%s/%v/%v\n", a.Lemma, a.Case, a.Number)
    }
}
```

## Склонение ФИО

`DeclineFullName` склоняет три компонента (фамилия / имя / отчество)
независимо. Любое из полей может быть пустым — пустые остаются пустыми
в результате. Род определяется по приоритету источников:
отчество → имя → фамилия. Для отчеств и фамилий, отсутствующих в словаре,
работает эвристика по окончанию.

```go
// Полное мужское ФИО.
out, _ := gomorphy.DeclineFullName(
    gomorphy.FullName{Last: "Иванов", First: "Иван", Patronymic: "Иванович"},
    gomorphy.Genitive,
)
fmt.Println(out.Last, out.First, out.Patronymic)
// Иванова Ивана Ивановича

// Полное женское ФИО (отчество не в словаре — работает эвристика).
out, _ = gomorphy.DeclineFullName(
    gomorphy.FullName{Last: "Иванова", First: "Анна", Patronymic: "Сергеевна"},
    gomorphy.Dative,
)
fmt.Println(out.Last, out.First, out.Patronymic)
// Ивановой Анне Сергеевне

// Несклоняемая иностранная фамилия + русское имя.
out, _ = gomorphy.DeclineFullName(
    gomorphy.FullName{Last: "Дюма", First: "Александр"},
    gomorphy.Instrumental,
)
fmt.Println(out.Last, out.First)
// Дюма Александром
```

## Разбор ФИО (любой падеж → именительный)

`ParseFullName` принимает строку «Last First Patronymic» (русский порядок),
«First Patronymic Last» (западный) или сокращённые комбинации, разбивает
её на компоненты и приводит каждый к им.п. — независимо от исходного падежа.
В паре с `DeclineFullName` это даёт стандартный кейс из реальных БД:
«ФИО в любом падеже → ФИО в любом другом падеже».

```go
// Из дательного падежа сразу получаем им.п.
nom, _ := gomorphy.ParseFullName("Ивановой Анне Сергеевне")
fmt.Println(nom.Last, nom.First, nom.Patronymic)
// Иванова Анна Сергеевна

// И склоняем в любой нужный.
abl, _ := gomorphy.DeclineFullName(nom, gomorphy.Instrumental)
fmt.Println(abl.Last, abl.First, abl.Patronymic)
// Ивановой Анной Сергеевной
```

Если структура уже разбита по полям — используйте `ToNominative`:

```go
nom, _ := gomorphy.ToNominative(gomorphy.FullName{
    Last: "Достоевским", First: "Фёдором", Patronymic: "Михайловичем",
})
// nom == FullName{Last: "Достоевский", First: "Фёдор", Patronymic: "Михайлович"}
```

Обратная эвристика умеет распознавать падежные формы без обращения к словарю
(притяжательные `-ов(а)/-ев(а)/-ин(а)`, адъективные `-ский/-цкий/-ой/-ая`,
несклоняемые `-ых/-их`, отчества всех типов). Двусмысленные формы
(«Иванова» — F Nom или M Gen?) разрешаются по подсказке из соседних
компонентов; без подсказки выбирается F (более частая трактовка).

## Лицензия

Проект распространяется по составной модели — разные части дистрибутива
покрыты разными лицензиями:

- **Исходный код и тестовые фикстуры** (всё под `*.go`, `cmd/`,
  `internal/`, `testdata/`): [MIT License](LICENSE).
  `testdata/sample.xml` — это ручная микровыборка, составленная с нуля:
  использован XML-формат OpenCorpora, но сами леммы и их формы —
  общеизвестные факты русского языка, не извлечены из дампа OpenCorpora.
- **Словарь OpenCorpora** (скачивается пользователем самостоятельно с
  [opencorpora.org](https://opencorpora.org/); `cmd/builddict` только
  конвертирует локальный XML-дамп в компактный `.bin`):
  [CC BY-SA](https://creativecommons.org/licenses/by-sa/4.0/),
  © OpenCorpora contributors. В репозиторий не включён.

Если вы поставляете собранный `.bin` (произведённый из дампа OpenCorpora)
в составе своего проекта, этот артефакт остаётся под CC BY-SA — MIT-лицензия
покрывает только исходный код, который такой артефакт создаёт и читает.

## Документация

Подробное описание формата словаря, обратного индекса и алгоритмов см. в
[docs/DESIGN.md](docs/DESIGN.md).
