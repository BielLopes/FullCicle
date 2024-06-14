import {MyWallet} from "../components/MyWallet";

export default async function HomePage({ params }: { params: { wallet_id: string } }) {
  return (
    <div>
      <h1>Meus Investimentos</h1>
      <p>{params.wallet_id}</p>
      <MyWallet wallet_id={params.wallet_id} />
    </div>
  )
}