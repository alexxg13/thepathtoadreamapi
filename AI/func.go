package AI

import (
	"fmt"
	"github.com/coalaura/mistral"
	"log"
)

func GenerationPromt(goal string, arg []string) string {
	return fmt.Sprintf(
		"Я хочу построить roadmap для достижения моих целей. Мои цели следующие: %s. Мне нужно, чтобы roadmap был подробным и включал как долгосрочные, так и краткосрочные этапы. Ответь на следующие вопросы, чтобы помочь мне сделать roadmap более точным и практичным:\n"+
			"1. Какая твоя главная цель? — Опиши самую важную цель, которая стоит перед тобой.\n"+
			"Ответ на вопрос: %s \n"+
			"2. Когда ты хочешь достичь этой цели? — В какой срок ты хочешь достичь этой цели (краткосрочный, среднесрочный, долгосрочный)?\n"+
			"Ответ на вопрос: %s \n"+
			"3. Что тебе нужно узнать или изучить для достижения этой цели? — Это поможет нам разбить roadmap на этапы.\n"+
			"Ответ на вопрос: %s \n"+
			"4. Какие ресурсы у тебя уже есть? — Например, знания, связи, навыки, финансы.\n"+
			"Ответ на вопрос: %s \n"+
			"5. Какие могут быть препятствия на пути? — Это поможет нам предусмотреть решения заранее.\n"+
			"Ответ на вопрос: %s \n"+
			"6. Какие промежуточные цели ты можешь поставить на пути к главной цели? — Мы разобьем цель на маленькие шаги.\n"+
			"Ответ на вопрос: %s \n"+
			"7. Как ты планируешь отслеживать свой прогресс? — Будем ли мы ставить контрольные точки на каждом этапе?\n"+
			"Ответ на вопрос: %s \n"+
			"8. Какие есть подцели или второстепенные задачи, которые помогут достичь главной цели?\n"+
			"Ответ на вопрос: %s \n",
		goal, arg[0], arg[1], arg[2], arg[3], arg[4], arg[5], arg[6], arg[7])
}

func SendPrompt(client *mistral.MistralClient, goal string, arg []string) (string, error) {
	prompt := GenerationPromt(goal, arg)

	request := mistral.ChatCompletionRequest{
		Model: mistral.ModelMistralSmall,
		Messages: []mistral.Message{
			{
				Role:    mistral.RoleSystem,
				Content: "Ты крутой наставник в сфере продвижения личного бренда и достижения любых целей",
			},
			{
				Role:    mistral.RoleUser,
				Content: prompt,
			},
		},
	}
	responce, err := client.Chat(request)
	if err != nil {
		log.Printf("Error responce Mistral AI: %s", err)
		return "", err
	}

	return responce.Choices[0].Message.Content, nil
}
