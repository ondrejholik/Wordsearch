package main

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


func random_letter() rune{
  var chars = []rune("abcdefghijklmnopqrstuvwxyz")
  return chars[rand.Intn(len(chars))]
}

func main(){
  args := os.Args
  var words []string
  var letter_count = 0
  var max_word = 0
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
    letter_count += len(text)
    if(len(text) > max_word){
      max_word = len(text)
    }
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }

  var size = max_word + 3
  tmp_size := int(math.Sqrt(float64(letter_count * 3)))
  if max_word < tmp_size {
    size = tmp_size
  }
  words_count := len(words)
  rand.Seed(time.Now().UnixNano())

  counter := 0
  matrix := make([][]rune, size)
  for i := range matrix{
    matrix[i] = make([]rune, size)
  }

  for i := 0; i < words_count; i++ {
    // Two random numbers from 0 to (size-1)
    counter = 0
    for  {
      x := rand.Intn(size)
      y := rand.Intn(size)
      counter++

      // Random number 0 - 7 => all direction including diagonals
      dir_num := rand.Intn(7)

      direction := directions[dir_num];
      tmp_x := x + direction[0] * len(words[i])
      tmp_y := y + direction[1] * len(words[i])
      exit := false
      if tmp_x >= 0 && tmp_x < size  && tmp_y >= 0 && tmp_y < size {
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


  var S = 165 * size
  dc := gg.NewContext(S, S)


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
     y := S/2 - h*len(matrix)/2 + i*h
     dc.DrawStringAnchored(line, float64(S/2), float64(y), 0.5, 0.5)
   }
   dc.SavePNG("wordsearch.png")
}
