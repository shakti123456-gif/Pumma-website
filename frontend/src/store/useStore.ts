import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Product } from '../data/products'
type CartItem={product:Product;size:string;quantity:number}
type Store={cart:CartItem[]; wishlist:string[]; cartOpen:boolean; addToCart:(product:Product,size?:string)=>void; removeFromCart:(id:string)=>void; toggleWish:(id:string)=>void; setCartOpen:(open:boolean)=>void}
export const useStore=create<Store>()(persist((set)=>({cart:[],wishlist:[],cartOpen:false,addToCart:(product,size='M')=>set(s=>{const item=s.cart.find(x=>x.product.id===product.id&&x.size===size);return {cart:item?s.cart.map(x=>x===item?{...x,quantity:x.quantity+1}:x):[...s.cart,{product,size,quantity:1}],cartOpen:true}}),removeFromCart:(id)=>set(s=>({cart:s.cart.filter(x=>x.product.id!==id)})),toggleWish:(id)=>set(s=>({wishlist:s.wishlist.includes(id)?s.wishlist.filter(x=>x!==id):[...s.wishlist,id]})),setCartOpen:(cartOpen)=>set({cartOpen})}),{name:'aero-store'}))
