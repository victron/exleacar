package paginator

var ALLOWED_DOMAINS = []string{"www.exleasingcar.com", "exleasingcar.com"}

const HOST = "exleasingcar.com"
const DETAILS_PREFIX = "https://www.exleasingcar.com/en/auto-details/"
const START_URL = "https://www.exleasingcar.com/en/auto-auction/order-9/"
const USER_AGENT = `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.77 Safari/537.36`
const CACHE_DIR = ".colly_cache2"
const DATA_DIR = "/home/vic/exle"
const MAX_PHOTOS_NUMBER = 20 // save no more then n photos

// mongo consts
const MONGO_LOCAL = "mongodb://localhost:27017"
const DB = "exlea"
const EXLE_CARS = "cars"

// sleep timers
const NEXT_PAGE = 10

// const NEXT_PAGE_RAND = 5
// const NEXT_CATEGORY = 30

const AUCTION_BLOCK = "div.auction-list-block"
