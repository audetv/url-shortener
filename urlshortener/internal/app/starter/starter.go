package starter

import (
	"context"
	"sync"

	"github.com/audetv/url-shortener/urlshortener/internal/app/repos/link"
)

// стартер приложение, оно не должно знать внешнее апи приложения - слой внешнего адаптера.
// Стартер - слой бизнес логики. Поэтому тоже делаем интерфейс HTTPServer interface и передаем его в Serve,
// точнее стартер - апликейшен уровень над уровнем бизнес логики.
// Здесь может располагаться бизнес логика по оркестрации запросов дополнительно.

// App Здесь мы должны стартануть приложение
type App struct {
	links *link.Links
}

// NewApp функция инициализации приложения, котора возвращает уже заполненный апп
// не хватает стора, получим его снаружи, пробросим в параметр user.UserStore
func NewApp(ls link.LinkStoreInterface) *App {
	a := &App{
		links: link.NewLinks(ls),
	}
	return a
}

type HTTPServer interface {
	Start(links *link.Links)
	Stop()
}

// Serve пока ничего не делает, ему на вход не хватает api, который будет присылать запросы
// Пробрасываем контекст, чтобы отловить сигналы от операционной системы
// и вейт группу, и в дефере завершаем ее
func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs HTTPServer) {
	defer wg.Done()
	// вызываем старт
	hs.Start(a.links)
	// дожидаемся здесь цтикс дана
	<-ctx.Done()
	// после того как все завершили, мы его останавливаем.
	// стоп добавит 2 сек и нормально остановит с бэкграунд контекстом,
	// уберем контекст, так как мы ничего не логируем, перенесем его в стоп
	hs.Stop()
}
