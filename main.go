package main

import "fmt"

func main() {

	// コンテンツのキャラクターのセリフデータ
	posts := []Post{
		{
			Name: "ケンシロウ",
			Text: "お前も負けた後のことを考えなかったようだな",
			Tags: map[string]string{"形式": "漫画", "作品": "北斗の拳"},
		},
		{
			Name: "ケンシロウ",
			Text: "オレの墓標に名はいらぬ！！死するならば戦いの荒野で！！",
			Tags: map[string]string{"形式": "漫画", "作品": "北斗の拳"},
		},
		{
			Name: "ケンシロウ",
			Text: "これからはまともに働いて食うんだな。それだけの筋力は残してある。",
			Tags: map[string]string{"形式": "漫画", "作品": "北斗の拳"},
		},
		{
			Name: "サウザー",
			Text: "愛ゆえに人は苦しまねばならぬ！！ 愛ゆえに人は悲しまねばならぬ！！ 愛ゆえに・・・",
			Tags: map[string]string{"形式": "漫画", "作品": "北斗の拳"},
		},
		{
			Name: "サウザー",
			Text: "ひとーーーーっふたーーーーっみいーーーーっ！！",
			Tags: map[string]string{"形式": "漫画", "作品": "北斗の拳"},
		},
		{
			Name: "サウザー",
			Text: "デカくなったな。小僧",
			Tags: map[string]string{"形式": "漫画", "作品": "北斗の拳"},
		},
		{
			Name: "タイラー・ダーデン",
			Text: "自分の殻を破るには、自分をまず壊すしかないようだ",
			Tags: map[string]string{"形式": "映画", "作品": "ファイトクラブ"},
		},
		{
			Name: "タイラー・ダーデン",
			Text: "ウジ虫どもうぬぼれるな！お前らは美しくもなければ特別なもんでもない。他と同様朽ち果てて消えるだけの有機物質だ",
			Tags: map[string]string{"形式": "映画", "作品": "ファイトクラブ"},
		},
		{
			Name: "タイラー・ダーデン",
			Text: "ルールその１、ファイト・クラブについて口にしてはならない",
			Tags: map[string]string{"形式": "映画", "作品": "ファイトクラブ"},
		},
		{
			Name: "碇シンジ",
			Text: "わかってる。内臓電源終了までの62秒でけりをつける",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "碇シンジ",
			Text: "…ちくしょう。ちくしょう…。ちくしょう ちくしょう",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "碇シンジ",
			Text: "裏切ったな！僕の気持ちを裏切ったな！父さんと同じに裏切ったんだ！",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "アスカ",
			Text: "あんたバカ？",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "アスカ",
			Text: "コネメガネ！いつまで歌ってんのよ鬱陶しい！",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "アスカ",
			Text: "逃げんなゴラァ！",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "碇ゲンドウ",
			Text: "ダメな時はレイを使うまでだ",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "碇ゲンドウ",
			Text: "時間がない。ATフィールドがお前の形を保てなくなる",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
		{
			Name: "碇ゲンドウ",
			Text: "では出て行け",
			Tags: map[string]string{"形式": "アニメ", "作品": "エヴァンゲリオン"},
		},
	}

	// 上のデータをフィルターする
	// 下がフィルタの定義

	// 作品=北斗の拳だけを抽出
	// filterGroup, err := Parse(`hasTag("作品", "北斗の拳")`)
	// 作品=北斗の拳、またはエヴァンゲリオンを抽出
	// filterGroup, err := Parse(`hasTag("作品", "北斗の拳") || hasTag("作品", "エヴァンゲリオン") `)
	// セリフに「な」を含む、または、名前に「ん」を含む
	// filterGroup, err := Parse(`includeWordInText("な") || includeWordInName("ン") `)
	// セリフに「な」を含む、かつ、名前に「ん」を含む
	// filterGroup, err := Parse(`includeWordInText("な") && includeWordInName("ン") `)
	// 複雑なフィルタもお手のもの
	filterGroup, err := Parse(`
	    (
			hasTag("作品", "エヴァンゲリオン") && (
				includeWordInName("碇シンジ")
				||
				includeWordInName("アスカ")
				&&
				(
					includeWordInText("！")
				)
			)
		)
		||
		(
			includeWordInName("タイラー・ダーデン")
		)
		||
		(
			hasTag("作品", "北斗の拳") && includeWordInText("っ")
		)
	`)
	if err != nil {
		panic(err)
	}

	filteredPosts := []Post{}
	for _, post := range posts {
		result, err := filterGroup.Eval(&post)
		if err != nil {
			panic(err)
		}
		if !result {
			continue
		}
		filteredPosts = append(filteredPosts, post)
	}

	// 結果
	printPosts(filteredPosts)
}

func printPosts(posts []Post) {
	for _, post := range posts {
		fmt.Printf("%s %s %s\n", post.Name, post.Text, post.Tags)
	}
}
