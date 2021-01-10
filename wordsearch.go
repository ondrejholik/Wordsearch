package wordsearch

import ( "fmt"
         "time"
         "os"
         "log"
         "bufio"
         "math"
         "math/rand"
         "github.com/fogleman/gg"
       )

/*
  get multiple words from user
  put words randomly in matrix
  print matrix
  export matrix with words to png 

*/


func randomLetter() rune{
  var chars = []rune("abcdefghijklmnopqrstuvwxyz")
  return chars[rand.Intn(len(chars))]
}

func main(){
  args := os.Args
  var words []string
  var letterCount = 0
  var maxWord = 0
  directions := [8][2]int{
    {-1,-1},
    {-1,0},
    {-1,1},
    {0,-1},
    {0,1},
    {1,-1},
    {1,0},
    {1,1},
  }

  if len(args) != 2{
    log.Println("Error. Bad number of arguments. One argument required(PATH)")
    os.Exit(1)
  }

  f, err := os.Open(args[1])

  if err != nil {
      log.Fatal(err)
  }

  defer f.Close()

  scanner := bufio.NewScanner(f)

  for scanner.Scan() {
    text := scanner.Text()
    words = append(words, text)
    letterCount += len(text)
    if(len(text) > maxWord){
      maxWord = len(text)
    }
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  var size = maxWord + 3
  tmpSize := int(math.Sqrt(float64(letterCount * 3)))
  if maxWord < tmpSize {
    size = tmpSize
  }
  wordsCount := len(words)
  rand.Seed(time.Now().UnixNano())

  counter := 0
  matrix := make([][]rune, size)
  for i := range matrix{
    matrix[i] = make([]rune, size)
  }

  for i := 0; i < wordsCount; i++ {
    // Two random numbers from 0 to (size-1)
    counter = 0
    for  {
      x := rand.Intn(size)
      y := rand.Intn(size)
      counter++

      // Random number 0 - 7 => all direction including diagonals
      dirNum := rand.Intn(7)

      direction := directions[dirNum];
      tmpX := x + direction[0] * len(words[i])
      tmpY := y + direction[1] * len(words[i])
      exit := false
      if tmpX >= 0 && tmpX < size  && tmpY >= 0 && tmpY < size {
        for c := 0; c < len(words[i]); c++{
          if matrix[x + (c*direction[0])][y + (c*direction[1])] != 0 &&
           !(matrix[x + (c*direction[0])][y + (c*direction[1])] == []rune(words[i])[c]){
            exit = true
            continue
          }
        }
        if !exit {
          for c := 0; c < len(words[i]); c++{
            matrix[x + (c*direction[0])][y + (c*direction[1])] = []rune(words[i])[c]
          }
        }else{
          continue
        }
        break
      }
      if counter > 10000{
        fmt.Println("Not solvable, Abort")
        break
        os.Exit(2)
      }
    }
  }

  // Fill empty space with random letters
  for i := 0; i < size; i++ {
    for j := 0; j < size; j++ {
      if matrix[i][j] == 0 {
        matrix[i][j] = random_letter()
      }
  }
}


  for i := 0; i < size; i++ {
    for j := 0; j < size; j++ {
        fmt.Print(string(matrix[i][j]))
    }
    fmt.Println()
  }


  var W = 165 * size
  var H = 165 * size + (size/4)*165
  dc := gg.NewContext(W, H)


   dc.SetRGB(1, 1, 1)
   dc.Clear()
   dc.SetRGB(0, 0, 0)
   if err := dc.LoadFontFace("./roboto.ttf", 128); err != nil {
     panic(err)
   }
   const h = 128
   for i := range matrix{
     line := ""
     for j := range matrix[i]{
       line += string(matrix[i][j])
     }
     y := W/2 - h*len(matrix)/2 + i*h
     dc.DrawStringAnchored(line, float64(W/2), float64(y), 0.5, 0.5)
   }

   const h2 = 64
   if err := dc.LoadFontFace("./roboto.ttf", 64); err != nil {
     panic(err)
   }
   var line string = ""
   var lc int = 0
   var wc int = 0
   for i := range words {
     if (wc + len(words[i])+2)*32 > int(float64(W)*0.8) {
        y := 165*size - (size/8)*165 + lc * h2
        dc.DrawStringAnchored(line, float64(W/2), float64(y), 0.5, 0.5)
        lc++
        wc = 0
        line = ""
     }
     line += words[i] + ", "
     wc += len(words[i])+2
   }
   y := 165*size - (size/8)*165 + lc*h2
   dc.DrawStringAnchored(line, float64(W/2), float64(y), 0.5, 0.5)
   dc.SavePNG("wordsearch.png")
}
