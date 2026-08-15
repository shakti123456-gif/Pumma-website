import type { Gender, ProductType, Sport } from './catalog'

export type Product = {
  id: string
  name: string
  sport: Sport
  gender: Gender
  type: ProductType
  price: number
  originalPrice?: number
  colors: string[]
  rating: number
  reviews: number
  image: string
  image2: string
  badge?: string
  description: string
  sizes: string[]
}

const pics = [
  ['photo-1552346154-21d32810aba3', 'photo-1542291026-7eec264c27ff'],
  ['photo-1608231387042-66d1773070a5', 'photo-1600185365483-26d7a4cc7519'],
  ['photo-1575537302964-96cd47c06b1b', 'photo-1600185365926-3a2ce3cdb9eb'],
  ['photo-1606107557195-0e29a4b5b4aa', 'photo-1603808033192-082d6919d3e1'],
  ['photo-1605348532760-6753d2c43329', 'photo-1460353581641-37baddab0fa2'],
  ['photo-1517836357463-d25dfeac3438', 'photo-1518611012118-696072aa579a'],
  ['photo-1538805060514-97d9cc17730c', 'photo-1523398002811-999ca8dec234'],
  ['photo-1515886657613-9f3515b0c78f', 'photo-1558618666-fcd25c85cd64'],
]

type Seed = { name: string; sport: Sport; gender: Gender; type: ProductType; price: number; sale?: boolean }

const seed: Seed[] = [
  // Running
  { name: 'Velocity NXT', sport: 'Running', gender: 'Men', type: 'Shoes', price: 125 },
  { name: 'RS-01 Runner', sport: 'Running', gender: 'Women', type: 'Shoes', price: 110, sale: true },
  { name: 'Sprint Elite', sport: 'Running', gender: 'Men', type: 'Clothing', price: 58 },
  { name: 'Cloud Knit Tee', sport: 'Running', gender: 'Women', type: 'Clothing', price: 45 },
  { name: 'Dash Runner Short', sport: 'Running', gender: 'Men', type: 'Clothing', price: 42 },
  { name: 'Arc Runner Legging', sport: 'Running', gender: 'Women', type: 'Clothing', price: 68 },
  // Training
  { name: 'Lift Training Tee', sport: 'Training', gender: 'Men', type: 'Clothing', price: 38 },
  { name: 'Reform Leggings', sport: 'Training', gender: 'Women', type: 'Clothing', price: 72 },
  { name: 'Core Half-Zip', sport: 'Training', gender: 'Men', type: 'Clothing', price: 65 },
  { name: 'Zone Sports Bra', sport: 'Training', gender: 'Women', type: 'Clothing', price: 48 },
  { name: 'Pro Training Bag', sport: 'Training', gender: 'Unisex', type: 'Clothing', price: 85 },
  { name: 'Studio Seamless Top', sport: 'Training', gender: 'Women', type: 'Clothing', price: 55, sale: true },
  // Football
  { name: 'Striker Pro Boot', sport: 'Football', gender: 'Men', type: 'Shoes', price: 140 },
  { name: 'Pitch Control Boot', sport: 'Football', gender: 'Women', type: 'Shoes', price: 130 },
  { name: 'Matchday Jersey', sport: 'Football', gender: 'Men', type: 'Clothing', price: 75 },
  { name: 'Elite Football Short', sport: 'Football', gender: 'Women', type: 'Clothing', price: 52 },
  // Basketball
  { name: 'Rally Court', sport: 'Basketball', gender: 'Men', type: 'Shoes', price: 115 },
  { name: 'Lumen Court', sport: 'Basketball', gender: 'Women', type: 'Shoes', price: 105 },
  { name: 'Hoops Mesh Short', sport: 'Basketball', gender: 'Men', type: 'Clothing', price: 48 },
  { name: 'Court Flex Tank', sport: 'Basketball', gender: 'Women', type: 'Clothing', price: 44 },
  // Lifestyle
  { name: 'Form Classic', sport: 'Lifestyle', gender: 'Men', type: 'Clothing', price: 78 },
  { name: 'Nova Luxe Hoodie', sport: 'Lifestyle', gender: 'Women', type: 'Clothing', price: 95 },
  { name: 'Everyday Cargo', sport: 'Lifestyle', gender: 'Men', type: 'Clothing', price: 88 },
  { name: 'Futuro Hoodie', sport: 'Lifestyle', gender: 'Women', type: 'Clothing', price: 82, sale: true },
  { name: 'Drift Suede', sport: 'Lifestyle', gender: 'Men', type: 'Shoes', price: 120 },
  { name: 'Vibe Slipstream', sport: 'Lifestyle', gender: 'Women', type: 'Shoes', price: 98 },
  { name: 'Orbit Windbreaker', sport: 'Lifestyle', gender: 'Men', type: 'Clothing', price: 140 },
  { name: 'Pulse Rib Top', sport: 'Lifestyle', gender: 'Women', type: 'Clothing', price: 62 },
]

export const products: Product[] = seed.map((item, i) => {
  const pair = pics[i % pics.length]
  return {
    id: `p${i + 1}`,
    name: item.name,
    sport: item.sport,
    gender: item.gender,
    type: item.type,
    price: item.price,
    originalPrice: item.sale ? item.price + 35 : undefined,
    colors: i % 2 ? ['#111111', '#b7b0a5', '#d5e6e9'] : ['#eef0ed', '#111111', '#b8694f'],
    rating: 4.4 + (i % 6) / 10,
    reviews: 12 + i * 7,
    image: `https://images.unsplash.com/${pair[0]}?auto=format&fit=crop&w=900&q=85`,
    image2: `https://images.unsplash.com/${pair[1]}?auto=format&fit=crop&w=900&q=85`,
    badge: item.sale ? '-30%' : i % 5 === 0 ? 'NEW' : undefined,
    description: `${item.sport} ${item.type.toLowerCase()} built for ${item.gender.toLowerCase()} athletes — lightweight, responsive, and ready for every session.`,
    sizes: item.type === 'Shoes' ? ['7', '8', '9', '10', '11'] : ['XS', 'S', 'M', 'L', 'XL'],
  }
})

export const featured = products.slice(0, 8)
