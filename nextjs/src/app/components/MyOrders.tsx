import { Order } from "../models";
//Server Components

async function getOrders(wallet_id: string): Promise<Order[]> {
  const response = await fetch(`http://localhost:8000/wallets/${wallet_id}/orders`);
  return response.json();
}

export async function MyOrders(props: { wallet_id: string } ) {

  const Orders = await getOrders(props.wallet_id);

  return (
    <ul>
      {Orders.map((order) => (
        <li key={order.id}>
          {order.Asset.id} - {order.shares} - R${" "} {order.price} - {order.status}
        </li>
      ))}
    </ul>
  );
}