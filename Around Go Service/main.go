  1 package main
  2 
  3 import (
  4 >       "cloud.google.com/go/storage"
  5 >       "context"
  6 >       "encoding/json"
  7 >       "fmt"
  8 >       jwtmiddleware "github.com/auth0/go-jwt-middleware"
  9 >       jwt "github.com/dgrijalva/jwt-go"
 10 >       "github.com/gorilla/mux"
 11 >       "github.com/olivere/elastic"
 12 >       "github.com/pborman/uuid"
 13 >       "io"
 14 >       "log"
 15 >       "net/http"
 16 >       "path/filepath"
 17 >       "reflect"
 18 >       "strconv"
 19 )
 20 
 21 const (
 22 >       POST_INDEX = "post"
 23 >       DISTANCE   = "200km"
 24 
 25 >       ES_URL      = "http://10.128.0.2:9200"
 26 >       BUCKET_NAME = "jz-bucket-132"
 27 )
 28 
 29 //for front end to render
 30 var (
 31 >       mediaTypes = map[string]string{
 32 >       >       ".jpeg": "image",
 33 >       >       ".jpg":  "image",
 34 >       >       ".gif":  "image",
 35 >       >       ".png":  "image",
 36 >       >       ".mov":  "video",
 37 >       >       ".mp4":  "video",
 38 >       >       ".avi":  "video",
 39 >       >       ".flv":  "video",
 40 >       >       ".wmv":  "video",
 41 >       }
 42 )
 43 
 44 type Post struct {
 45 >       User     string   `json:"user"`
 46 >       Message  string   `json:"message"`
 47 >       Url      string   `json:"url"`
 48 >       Type     string   `json:"type"`
 49 >       Face     float32  `json:"face"`
 50 >       Location Location `json:"location"`
 51 }
 52 
 53 type Location struct {
 54 >       Lat float64 `json:"lat"`
 55 >       Lon float64 `json:"lon"`
 56 }
 57 
 58 func handlerPost(w http.ResponseWriter, r *http.Request) {
 59 >       // Parse from body of request to get a json object
 60 >       fmt.Println("Received one post request")
 61 >       w.Header().Set("Content-Type", "application/json")
 62 >       w.Header().Set("Access-Control-Allow-Origin", "*")
 63 >       w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
 64 
 65 >       if r.Method == "OPTIONS" {
 66 >       >       return
 67 >       }
 68 
 69 >       //user token
 70 >       token := r.Context().Value("user")
 71 >       claims := token.(*jwt.Token).Claims
 72 >       username := claims.(jwt.MapClaims)["username"]
 73 
 74 >       //read parameters from request
 75 >       lat, _ := strconv.ParseFloat(r.FormValue("lat"), 64)
 76 >       lon, _ := strconv.ParseFloat(r.FormValue("lon"), 64)
 77 
 78 >       p := &Post{
 79 >       >       User:    username.(string),
 80 >       >       Message: r.FormValue("message"),
 81 >       >       Location: Location{
 82 >       >       >       Lat: lat,
 83 >       >       >       Lon: lon,
 84 >       >       },
 85 >       >       Face: 0.0,
 86 >       }
 87 
 88 >       file, header, err := r.FormFile("image")
 89 >       if err != nil {
 90 >       >       http.Error(w, "Image is not available", http.StatusBadRequest)
 91 >       >       fmt.Printf("Image is not available %v\n", err)
 92 >       >       return
 93 >       }
 94 
 95 >       suffix := filepath.Ext(header.Filename)
 96 >       if t, ok := mediaTypes[suffix]; ok {
 97 >       >       p.Type = t
 98 >       } else {
 99 >       >       p.Type = "unknown"
100 >       }
101 
102 >       id := uuid.New()
103 >       attrs, err := saveToGCS(file, id)
104 >       if err != nil {
105 >       >       http.Error(w, "Failed to save iamge to GCS", http.StatusInternalServerError)
106 >       >       fmt.Printf("Failed to save image to GCS %v\n", err)
107 >       >       return
108 >       }
109 
110 >       p.Url = attrs.MediaLink
111 
112 >       if p.Type == "image" {
113 >       >       uri := fmt.Sprintf("gs://%s/%s", BUCKET_NAME, id)
114 >       >       if score, err := annotate(uri); err != nil {
115 >       >       >       http.Error(w, "Failed to annotate image", http.StatusInternalServerError)
116 >       >       >       fmt.Printf("Failed to annotate image %v\n", err)
117 >       >       >       return
118 >       >       } else {
119 >       >       >       p.Face = score
120 >       >       }
121 >       }
122 
123 >       err = saveToES(p, POST_INDEX, id)
124 >       if err != nil {
125 >       >       http.Error(w, "Failed to save post to Elasticsearch", http.StatusInternalServerError)
126 >       >       fmt.Printf("Failed to save post to Elasticsearch %v\n", err)
127 >       >       return
128 >       }
129 
130 >       fmt.Printf("Post is done successfully: %s\n", p.Message)
131 }
132 
133 func handlerSearch(w http.ResponseWriter, r *http.Request) {
134 >       fmt.Println("Received one search request")
135 
136 >       w.Header().Set("Content-Type", "application/json")
137 >       w.Header().Set("Access-Control-Allow-Origin", "*")
138 >       w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
139 
140 >       if r.Method == "OPTIONS" {
141 >       >       return
142 >       }
143 
144 >       // read parameters from request
145 >       lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
146 >       lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
147 
148 >       ran := DISTANCE
149 >       if val := r.URL.Query().Get("range"); val != "" {
150 >       >       ran = val + "km"
151 >       }
152 >       fmt.Println("Range is ", ran)
153 
154 >       // read post from Elasticsearch
155 >       query := elastic.NewGeoDistanceQuery("location")
156 >       query = query.Distance(ran).Lat(lat).Lon(lon)
157 
158 >       searchResult, err := readFromES(query, POST_INDEX)
159 >       if err != nil {
160 >       >       http.Error(w, "Failed to read from Elasticsearch", http.StatusInternalServerError)
161 >       >       fmt.Printf("Failed to read from Elasticsearch %v\n", err)
162 >       >       return
163 >       }
164 
165 >       posts := getPostFromSearchResult(searchResult)
166 
167 >       // convert read result to JSON and put in response
168 >       js, err := json.Marshal(posts)
169 >       if err != nil {
170 >       >       http.Error(w, "Failed to parse posts into json format", http.StatusInternalServerError)
171 >       >       fmt.Printf("Failed to parse posts into json format %v\n", err)
172 >       >       return
173 >       }
174 
175 >       w.Write(js)
176 }
177 
178 func handlerCluster(w http.ResponseWriter, r *http.Request) {
179 >       fmt.Println("Received one cluster request")
180 >       w.Header().Set("Content-Type", "application/json")
181 >       w.Header().Set("Access-Control-Allow-Origin", "*")
182 >       w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
183 
184 >       if r.Method == "OPTIONS" {
185 >       >       return
186 >       }
187 
188 >       term := r.URL.Query().Get("term")
189 >       //score >= 0.9
190 >       query := elastic.NewRangeQuery(term).Gte(0.9)
191 
192 >       searchResult, err := readFromES(query, POST_INDEX)
193 >       if err != nil {
194 >       >       http.Error(w, "Failed to read post from Elasticsearch", http.StatusInternalServerError)
195 >       >       fmt.Printf("Failed to read post from Elasticsearch %v.\n", err)
196 >       >       return
197 >       }
198 >       posts := getPostFromSearchResult(searchResult)
199 
200 >       js, err := json.Marshal(posts)
201 >       if err != nil {
202 >       >       http.Error(w, "Failed to parse post object", http.StatusInternalServerError)
203 >       >       fmt.Printf("Failed to parse post object %v\n", err)
204 >       >       return
205 >       }
206 
207 >       w.Write(js)
208 }
209 
210 func readFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
211 >       client, err := elastic.NewClient(elastic.SetURL(ES_URL))
212 >       if err != nil {
213 >       >       return nil, err
214 >       }
215 
216 >       searchResult, err := client.Search().
217 >       >       Index(index).
218 >       >       Query(query).
219 >       >       Pretty(true).
220 >       >       Do(context.Background())
221 >       if err != nil {
222 >       >       return nil, err
223 >       }
224 
225 >       return searchResult, nil
226 }
227 
228 func getPostFromSearchResult(searchResult *elastic.SearchResult) []Post {
229 >       var posts []Post
230 >       var ptype Post
231 
232 >       for _, item := range searchResult.Each(reflect.TypeOf(ptype)) {
233 >       >       if p, ok := item.(Post); ok {
234 >       >       >       posts = append(posts, p)
235 >       >       }
236 >       }
237 
238 >       return posts
239 }
240 
241 func saveToGCS(r io.Reader, objectName string) (*storage.ObjectAttrs, error) {
242 >       ctx := context.Background()
243 
244 >       client, err := storage.NewClient(ctx)
245 >       if err != nil {
246 >       >       return nil, err
247 >       }
248 
249 >       //bucket instance
250 >       bucket := client.Bucket(BUCKET_NAME)
251 >       //if cannot get bucket attributes => bucket does not exist
252 >       if _, err := bucket.Attrs(ctx); err != nil {
253 >       >       return nil, err
254 >       }
255 
256 >       //object instance
257 >       object := bucket.Object(objectName)
258 >       //upload file
259 >       wc := object.NewWriter(ctx)
260 >       if _, err := io.Copy(wc, r); err != nil {
261 >       >       return nil, err
262 >       }
263 
264 >       if err := wc.Close(); err != nil {
265 >       >       return nil, err
266 >       }
267 
268 >       //get attribute (url)  => ES
269 >       //uploaded file access => all users can read
270 >       if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
271 >       >       return nil, err
272 >       }
273 
274 >       attrs, err := object.Attrs(ctx)
275 >       if err != nil {
276 >       >       return nil, err
277 >       }
278 
279 >       fmt.Printf("Image is saved to GCS: %s\n", attrs.MediaLink)
280 >       return attrs, nil
281 }
282 
283 func saveToES(i interface{}, index string, id string) error {
284 >       client, err := elastic.NewClient(elastic.SetURL(ES_URL))
285 >       if err != nil {
286 >       >       return err
287 >       }
288 
289 >       _, err = client.Index().
290 >       >       Index(index).
291 >       >       Id(id).
292 >       >       BodyJson(i).
293 >       >       Do(context.Background())
294 
295 >       if err != nil {
296 >       >       return err
297 >       }
298 
299 >       return nil
300 }
301 
302 func main() {
303 >       fmt.Println("started-services")
304 
305 >       jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
306 >       >       ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
307 >       >       >       return []byte(mySigningKey), nil
308 >       >       },
309 >       >       SigningMethod: jwt.SigningMethodHS256,
310 >       })
311 
312 >       r := mux.NewRouter()
313 
314 >       r.Handle("/post", jwtMiddleware.Handler(http.HandlerFunc(handlerPost))).Methods("POST", "OPTIONS")
315 >       r.Handle("/search", jwtMiddleware.Handler(http.HandlerFunc(handlerSearch))).Methods("GET", "OPTIONS")
316 >       r.Handle("/cluster", jwtMiddleware.Handler(http.HandlerFunc(handlerCluster))).Methods("GET", "OPTIONS")
317 >       r.Handle("/signup", http.HandlerFunc(handlerSignup)).Methods("POST", "OPTIONS")
318 >       r.Handle("/login", http.HandlerFunc(handlerLogin)).Methods("POST", "OPTIONS")
319 
320 >       log.Fatal(http.ListenAndServe(":8080", r))
321 }

