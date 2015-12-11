package me.wener.bbvm.asm;

import com.google.common.base.MoreObjects;
import com.google.common.collect.Lists;
import me.wener.bbvm.vm.Operand;
import me.wener.bbvm.vm.Symbol;

import java.util.List;

/**
 * @author wener
 * @since 15/12/11
 */
public class Label implements Symbol {
    final String name;
    int address = -1;
    List<OperandInfo> operands = Lists.newArrayList();
    Token token;

    public Label(String name, Token token) {
        this.token = token;
        this.name = name;
    }

    public Label(String name) {
        this.name = name;
    }

    public int getAddress() {
        return address;
    }

    public Label setAddress(int address) {
        this.address = address;
        for (OperandInfo info : operands) {
            info.operand.setValue(address);
        }
        return this;
    }

    @Override
    public String getName() {
        return name;
    }

    @Override
    public String toString() {
        return MoreObjects.toStringHelper(this)
                .add("name", name)
                .add("address", address)
                .add("line", token != null ? token.beginLine : -1)
                .add("column", token != null ? token.beginColumn : -1)
                .toString();
    }

    public void addOperand(Token token, Operand operand) {
        operand.setSymbol(this).setValue(address);
        operands.add(new Label.OperandInfo(operand, token));
    }

    static class OperandInfo {
        Operand operand;
        Token token;

        public OperandInfo(Operand operand, Token token) {
            this.operand = operand;
            this.token = token;
        }
    }
}
