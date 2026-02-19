import 'package:intl/intl.dart';
import '../../data/models/transaction.dart';

class AmountUtils {
  static String formatAmount(double amount, {String currency = 'CNY'}) {
    final format = NumberFormat.currency(
      locale: 'zh_CN',
      symbol: currency == 'CNY' ? '¥' : '\$',
      decimalDigits: 2,
    );
    return format.format(amount);
  }

  static String formatWithSign(double amount, TransactionType type) {
    final formatted = formatAmount(amount);
    switch (type) {
      case TransactionType.income:
        return '+$formatted';
      case TransactionType.expense:
        return '-$formatted';
      case TransactionType.transfer:
        return formatted;
    }
  }

  static double roundToCents(double amount) {
    return (amount * 100).round() / 100;
  }

  static String formatAmountWithoutSymbol(double amount) {
    return NumberFormat('#,##0.00').format(amount);
  }
}
