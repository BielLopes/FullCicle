import { Module } from '@nestjs/common';
import { OrdersService } from './orders.service';
import { OrdersController } from './orders.controller';
import { ClientsModule, Transport } from '@nestjs/microservices';

@Module({
  imports: [
    // Configuração para publicar mensagens no kafka
    ClientsModule.register([
      {
        name: 'ORDERS_PUBLISHER',
        transport: Transport.KAFKA,
        options: {
          client: {
            clientId: 'orders',
            brokers: ['host.docker.internal:9094'],
          },
          producerOnlyMode: true,
        },
      },
    ]),
  ],

  controllers: [OrdersController],
  providers: [OrdersService],
})
export class OrdersModule {}
