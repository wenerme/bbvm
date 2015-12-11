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
public class Label extends AbstractAssembly implements Symbol, Assembly {
    final String name;
    int address = -1;
    List<OperandInfo> operands = Lists.newArrayList();
    private Token token;

    public Label(String name, Token token) {
        this.token = token;
        this.name = name;
    }

    public Label(String name) {
        this.name = name;
    }

    public Token getToken() {
        return token;
    }

    public Label setToken(Token token) {
        this.token = token;
        return this;
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
                .add("comment", comment)
                .add("line", token != null ? token.beginLine : -1)
                .add("column", token != null ? token.beginColumn : -1)
                .toString();
    }

    public void addOperand(Token token, Operand operand) {
        operand.setSymbol(this).setValue(address);
        operands.add(new Label.OperandInfo(operand, token));
    }

    @Override
    public Type getType() {
        return Type.LABEL;
    }

    @Override
    public String toAssembly() {
        return token.toString() + commentAssembly();
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
